package ioc

import (
	followv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/follow/v1"
	"gitee.com/geekbang/basic-go/webook/feed/repository"
	"gitee.com/geekbang/basic-go/webook/feed/service"
)

func RegisterHandler(repo repository.FeedEventRepo, followClient followv1.FollowServiceClient) map[string]service.Handler {
	articleHandler := service.NewArticleEventHandler(repo, followClient)
	followHanlder := service.NewFollowEventHandler(repo)
	likeHandler := service.NewLikeEventHandler(repo)
	return map[string]service.Handler{
		service.ArticleEventName: articleHandler,
		service.FollowEventName:  followHanlder,
		service.LikeEventName:    likeHandler,
	}
}
