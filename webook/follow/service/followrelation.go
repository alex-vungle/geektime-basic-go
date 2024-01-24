package service

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/follow/domain"
	"gitee.com/geekbang/basic-go/webook/follow/repository"
)

type FollowRelationService interface {
	GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
	FollowInfo(ctx context.Context,
		follower, followee int64) (domain.FollowRelation, error)
	Follow(ctx context.Context, follower, followee int64) error
	CancelFollow(ctx context.Context, follower, followee int64) error
}

type followRelationService struct {
	repo repository.FollowRepository
}

func (f *followRelationService) CancelFollow(ctx context.Context, follower, followee int64) error {
	return f.repo.InactiveFollowRelation(ctx, follower, followee)
}

func NewFollowRelationService(repo repository.FollowRepository) FollowRelationService {
	return &followRelationService{
		repo: repo,
	}
}

func (f *followRelationService) GetFollowee(ctx context.Context,
	follower, offset, limit int64) ([]domain.FollowRelation, error) {
	return f.repo.GetFollowee(ctx, follower, offset, limit)
}

func (f *followRelationService) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	val, err := f.repo.FollowInfo(ctx, follower, followee)
	return val, err
}

func (f *followRelationService) Follow(ctx context.Context, follower, followee int64) error {
	return f.repo.AddFollowRelation(ctx, domain.FollowRelation{
		Followee: followee,
		Follower: follower,
	})
}
