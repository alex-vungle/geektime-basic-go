package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"golang.org/x/sync/errgroup"
	"time"
)

type CachedGRPCArticleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache

	userRepo UserRepository
	intrRepo InteractiveRepository
}

func (c *CachedGRPCArticleRepository) ToDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime:  time.UnixMilli(art.Ctime),
		Utime:  time.UnixMilli(art.Utime),
		Status: domain.ArticleStatus(art.Status),
	}
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

func NewCachedGRPCArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache, userRepo UserRepository, intrRepo InteractiveRepository) *CachedGRPCArticleRepository {
	return &CachedGRPCArticleRepository{dao: dao, cache: cache, userRepo: userRepo, intrRepo: intrRepo}
}

func (c *CachedGRPCArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) Update(ctx context.Context, art domain.Article) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error {
	//TODO implement me
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CachedGRPCArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.ToDomain(art), nil
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

	// Write back to the cache
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		er := c.cache.SetPub(ctx, res)
		if er != nil {
			//Log an error
		}
	}()

	return res, nil
}

func (c *CachedGRPCArticleRepository) ListPub(ctx context.Context, uid int64, start time.Time, offset int, limit int) ([]domain.Article, error) {

	daoArts, err := c.dao.ListPub(ctx, start, offset, limit)

	if err != nil {
		return []domain.Article{}, err
	}
	domainArticles := make([]domain.Article, len(daoArts))

	for i, d := range daoArts {
		domainArticles[i] = c.ToDomain(dao.Article(d))
	}

	var (
		eg errgroup.Group
	)

	eg.Go(func() error {
		for _, art := range domainArticles {
			author, err1 := c.userRepo.FindById(ctx, art.Author.Id)
			if err1 != nil {
				// log an error
			}
			art.Author.Name = author.Nickname
		}
		return nil
	})

	eg.Go(func() error {
		for _, art := range domainArticles {
			intr, err1 := c.intrRepo.GetById(ctx, uid, art.Id)
			if err1 != nil {
				//LOg an error
			}
			art.Intr = intr
		}
		return nil
	})

	return domainArticles, nil
}
