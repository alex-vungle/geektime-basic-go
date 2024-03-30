package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"golang.org/x/sync/errgroup"
	"time"
)

// CachedGRPCArticleRepository 第十二周作业演示
type CachedGRPCArticleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache

	userRepo UserRepository
	intrRepo InteractiveRepository
}

func NewCachedGRPCArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, userRepo UserRepository, intrRepo InteractiveRepository) *CachedGRPCArticleRepository {
	return &CachedGRPCArticleRepository{dao: dao, cache: cache, userRepo: userRepo, intrRepo: intrRepo}
}

func (c *CachedGRPCArticleRepository) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]domain.Article, error) {
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) GetPubById(ctx context.Context, uid, id int64) (domain.Article, error) {
	res, err := c.cache.GetPub(ctx, id)
	if err == nil {
		return res, err
	}
	art, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 我现在要去查询对应的 User 信息，拿到创作者信息
	res = c.ToDomain(dao.Article(art))
	var (
		eg errgroup.Group
	)

	eg.Go(func() error {
		author, err1 := c.userRepo.FindById(ctx, art.AuthorId)
		res.Author.Name = author.Nickname
		return err1
	})

	eg.Go(func() error {
		intr, err1 := c.intrRepo.GetById(ctx, uid, id)
		res.Intr = intr
		return err1
	})

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := c.cache.SetPub(ctx, res)
		if er != nil {
			// 记录日志
		}
	}()
	return res, nil
}

func (c *CachedGRPCArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := c.cache.Get(ctx, id)
	if err == nil {
		return res, nil
	}
	art, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	res = c.ToDomain(art)
	go func() {
		er := c.cache.Set(ctx, res)
		if er != nil {
			// 记录日志
		}
	}()
	return res, nil
}

func (c *CachedGRPCArticleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error {
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) Update(ctx context.Context, art domain.Article) error {
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		//Status:   uint8(art.Status),
		Status: art.Status.ToUint8(),
	}
}

func (c *CachedGRPCArticleRepository) ToDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			// 这里有一个错误
			Id: art.AuthorId,
		},
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
		Status: domain.ArticleStatus(art.Status),
	}
}
