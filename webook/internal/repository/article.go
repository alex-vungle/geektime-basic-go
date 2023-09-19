package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error

	// Sync 本身要求先保存到制作库，再同步到线上库
	Sync(ctx context.Context, art domain.Article) (int64, error)
	// SyncStatus 仅仅同步状态
	SyncStatus(ctx context.Context, id int64, status domain.ArticleStatus) error
}

type CachedArticleRepository struct {
	// 操作单一的库
	dao dao.ArticleDAO

	// SyncV1 用
	authorDAO dao.ArticleAuthorDAO
	readerDAO dao.ArticleReaderDAO

	// SyncV2 用
	db *gorm.DB
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func NewArticleRepositoryV1(authorDAO dao.ArticleAuthorDAO,
	readerDAO dao.ArticleReaderDAO) ArticleRepository {
	return &CachedArticleRepository{
		authorDAO: authorDAO,
		readerDAO: readerDAO,
	}
}

func (repo *CachedArticleRepository) SyncStatus(ctx context.Context, id int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, id, status.ToUint8())
}

func (repo *CachedArticleRepository) Sync(ctx context.Context,
	art domain.Article) (int64, error) {
	return repo.dao.Sync(ctx, repo.toEntity(art))
}

func (repo *CachedArticleRepository) SyncV2(ctx context.Context,
	art domain.Article) (int64, error) {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 直接 defer Rollback
	// 如果我们后续 Commit 了，这里会得到一个错误，但是没关系
	defer tx.Rollback()
	authorDAO := dao.NewGORMArticleDAO(tx)
	readerDAO := dao.NewGORMArticleReaderDAO(tx)

	// 下面代码和 SyncV1 一模一样
	artn := repo.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = authorDAO.Create(ctx, artn)
		if err != nil {
			return 0, err
		}
	} else {
		err = authorDAO.UpdateById(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = readerDAO.UpsertV2(ctx, dao.PublishedArticle{Article: artn})
	if err != nil {
		// 依赖于 defer 来 rollback
		return 0, err
	}
	tx.Commit()
	return artn.Id, nil
}

func (repo *CachedArticleRepository) SyncV1(ctx context.Context,
	art domain.Article) (int64, error) {
	artn := repo.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = repo.authorDAO.Create(ctx, artn)
		if err != nil {
			return 0, err
		}
	} else {
		err = repo.authorDAO.UpdateById(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = repo.readerDAO.Upsert(ctx, artn)
	return id, err
}

func (repo *CachedArticleRepository) Create(ctx context.Context,
	art domain.Article) (int64, error) {
	return repo.dao.Create(ctx, repo.toEntity(art))
}

func (repo *CachedArticleRepository) Update(ctx context.Context,
	art domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func (repo *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		// 这一步，就是将领域状态转化为存储状态。
		// 这里我们就是直接转换，
		// 有些情况下，这里可能是借助一个 map 来转
		Status: uint8(art.Status),
	}
}
