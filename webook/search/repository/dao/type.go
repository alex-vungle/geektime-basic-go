package dao

import (
	"context"
)

type UserDAO interface {
	InputUser(ctx context.Context, user User) error
	Search(ctx context.Context, keywords []string) ([]User, error)
}

type ArticleDAO interface {
	InputArticle(ctx context.Context, article Article) error
	// Search artIds 命中了索引的 article id
	Search(ctx context.Context, req SearchReq, keywords []string) ([]Article, error)
}

type TagDAO interface {
	Search(ctx context.Context, uid int64, biz string, keywords []string) ([]int64, error)
}

type AnyDAO interface {
	Input(ctx context.Context, index, docID, data string) error
	Delete(ctx context.Context, index string, docID string) error
}

type LikeDAO interface {
	Search(ctx context.Context, uid int64, biz string) ([]int64, error)
}

type CollectDAO interface {
	Search(ctx context.Context, uid int64, biz string) ([]int64, error)
}

type SearchReq struct {
	LikeIds    []int64
	TagIds     []int64
	CollectIds []int64
}
