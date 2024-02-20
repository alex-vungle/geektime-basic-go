package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	events "gitee.com/geekbang/basic-go/webook/internal/events/article"
	"gitee.com/geekbang/basic-go/webook/internal/repository/article"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"time"
)

//go:generate mockgen -source=article.go -package=svcmocks -destination=mocks/article.mock.go ArticleService
type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	// ListPub 根据这个 start 时间来查询
	ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error)
}

type articleService struct {
	repo article.ArticleRepository

	// V1 依靠两个不同的 repository 来解决这种跨表，或者跨库的问题
	author   article.ArticleAuthorRepository
	reader   article.ArticleReaderRepository
	l        logger.LoggerV1
	producer events.Producer

	ch chan readInfo
}

func (svc *articleService) ListPub(ctx context.Context,
	start time.Time, offset, limit int) ([]domain.Article, error) {
	return svc.repo.ListPub(ctx, start, offset, limit)
}

type readInfo struct {
	uid int64
	aid int64
}

// GetPublishedByIdV1 批量发送的例子
func (svc *articleService) GetPublishedByIdV1(ctx context.Context, id, uid int64) (domain.Article, error) {
	// 另一个选项，在这里组装 Author，调用 UserService
	art, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			// 改批量的做法
			svc.ch <- readInfo{
				aid: id,
				uid: uid,
			}
		}()
	}
	return art, err
}

func (svc *articleService) batchSendReadInfo(ctx context.Context) {
	// 10 个一批
	// 单个转批量都要考虑的兜底问题
	for {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		const batchSize = 10
		uids := make([]int64, 0, 10)
		aids := make([]int64, 0, 10)
		send := false
		for !send {
			select {
			// 这边是超时了
			case <-ctx.Done():
				// 也要执行发送
				//goto send
				send = true
			case info, ok := <-svc.ch:
				if !ok {
					cancel()
					send = true
					continue
				}
				uids = append(uids, info.uid)
				aids = append(aids, info.aid)
				// 凑够了
				if len(uids) == batchSize {
					//goto send
					send = true
				}
			}
		}
		//send:
		// 装满了，凑够了一批
		svc.producer.ProduceReadEventV1(context.Background(),
			events.ReadEventV1{
				Uids: uids,
				Aids: aids,
			})
		cancel()
	}
}

func (svc *articleService) GetPublishedById(ctx context.Context, id, uid int64) (domain.Article, error) {
	// 另一个选项，在这里组装 Author，调用 UserService
	art, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		// 每次打开一篇文章，就发一条消息
		go func() {
			// 生产者也可以通过改批量来提高性能
			er := svc.producer.ProduceReadEvent(
				ctx,
				events.ReadEvent{
					// 即便你的消费者要用 art 的里面的数据，
					// 让它去查询，你不要在 event 里面带
					Uid: uid,
					Aid: id,
				})
			if er != nil {
				svc.l.Error("发送读者阅读事件失败")
			}
		}()

		//go func() {
		//	// 改批量的做法
		//	svc.ch <- readInfo{
		//		aid: id,
		//		uid: uid,
		//	}
		//}()
	}
	return art, err
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *articleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

func (a *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	// art.Status = domain.ArticleStatusPrivate 然后直接把整个 art 往下传
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	// 制作库
	//id, err := a.repo.Create(ctx, art)
	//// 线上库呢？
	//a.repo.SyncToLiveDB(ctx, art)
	return a.repo.Sync(ctx, art)
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("部分失败，保存到线上库失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}
	if err != nil {
		a.l.Error("部分失败，重试彻底失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
		// 接入你的告警系统，手工处理一下
		// 走异步，我直接保存到本地文件
		// 走 Canal
		// 打 MQ
	}
	return id, err
}

func NewArticleService(repo article.ArticleRepository,
	l logger.LoggerV1,
	producer events.Producer) ArticleService {
	res := &articleService{
		repo:     repo,
		producer: producer,
		l:        l,
		//ch:       make(chan readInfo, 10),
	}
	return res
}

// ctx, cancel := context.WithCancel(context.Background())
// NewArticleServiceV3(ctx)
// 这里一大堆业务逻辑
// 主程序（main 函数准备退出）
// cancel()
func NewArticleServiceV3(ctx context.Context, repo article.ArticleRepository,
	l logger.LoggerV1,
	producer events.Producer) ArticleService {
	res := &articleService{
		repo:     repo,
		producer: producer,
		l:        l,
		//ch:       make(chan readInfo, 10),
	}
	go func() {
		// 我系统关闭的时候，你 channel 里面还有数据，没发出去，怎么办？
		// 第一种是啥也不干，你在关闭的时候，time.Sleep
		// 第二种
		res.batchSendReadInfo(ctx)
	}()
	return res
}

func (a *articleService) Close() error {
	close(a.ch)
	return nil
}

func NewArticleServiceV2(repo article.ArticleRepository,
	l logger.LoggerV1,
	producer events.Producer) ArticleService {
	ch := make(chan readInfo, 10)
	go func() {
		for {
			uids := make([]int64, 0, 10)
			aids := make([]int64, 0, 10)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			for i := 0; i < 10; i++ {
				select {
				case info, ok := <-ch:
					if !ok {
						cancel()
						return
					}
					uids = append(uids, info.uid)
					aids = append(aids, info.aid)
				case <-ctx.Done():
					break
				}
			}
			cancel()
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			producer.ProduceReadEventV1(ctx, events.ReadEventV1{
				Uids: uids,
				Aids: aids,
			})
			cancel()
		}
	}()
	return &articleService{
		repo:     repo,
		producer: producer,
		l:        l,
		ch:       ch,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository,
	reader article.ArticleReaderRepository, l logger.LoggerV1) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
		l:      l,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a *articleService) update(ctx context.Context, art domain.Article) error {
	// 只要你不更新 author_id
	// 但是性能比较差
	//artInDB := a.repo.FindById(ctx, art.Id)
	//if art.Author.Id != artInDB.Author.Id {
	//	return errors.New("更新别人的数据")
	//}
	return a.repo.Update(ctx, art)
}
