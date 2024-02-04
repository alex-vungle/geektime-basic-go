package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/search/domain"
	"gitee.com/geekbang/basic-go/webook/search/repository"
	"golang.org/x/sync/errgroup"
	"strings"
)

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type searchService struct {
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
}

func NewSearchService(userRepo repository.UserRepository, articleRepo repository.ArticleRepository) SearchService {
	return &searchService{userRepo: userRepo, articleRepo: articleRepo}
}

func (s *searchService) Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error) {
	// 你要搜索用户，你也要搜索 article
	// 要对 expression 进行解析，生成查询计划
	// 输入预处理
	// 清除掉空格，切割;',.
	keywords := strings.Split(expression, " ")
	var eg errgroup.Group
	var res domain.SearchResult
	eg.Go(func() error {
		users, err := s.userRepo.SearchUser(ctx, keywords)
		res.Users = users
		return err
	})
	eg.Go(func() error {
		arts, err := s.articleRepo.SearchArticle(ctx, uid, keywords)
		res.Articles = arts
		return err
	})
	return res, eg.Wait()
}
