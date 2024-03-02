package service

import (
	"context"
	"fmt"
	followv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/follow/v1"
	"gitee.com/geekbang/basic-go/webook/feed/domain"
	"gitee.com/geekbang/basic-go/webook/feed/repository"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
)

type feedService struct {
	repo repository.FeedEventRepo
	// 对应的 string 就是 type
	handlerMap   map[string]Handler
	followClient followv1.FollowServiceClient
}

func NewFeedService(repo repository.FeedEventRepo, handlerMap map[string]Handler) FeedService {
	return &feedService{
		repo:       repo,
		handlerMap: handlerMap,
	}
}

func (f *feedService) RegisterService(typ string, handler Handler) {
	f.handlerMap[typ] = handler
}

//func (f *feedService) CreateLikeFeedEvent() {
//
//}

func (f *feedService) CreateFeedEvent(ctx context.Context, feed domain.FeedEvent) error {
	handler, ok := f.handlerMap[feed.Type]
	if !ok {
		// 说明 type 不对
		// 你还可以考虑兜底机制
		// 有一个 defaultHandler，然后调用 defaultHandler
		return fmt.Errorf("未能找到对应的 Handler %s", feed.Type)
	}
	return handler.CreateFeedEvent(ctx, feed.Ext)
}

// GetFeedEventListV1 不依赖于 Handler 的直接查询
func (f *feedService) GetFeedEventListV1(ctx context.Context, uid int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	// 直接查询
	var eg errgroup.Group
	var lock sync.Mutex
	events := make([]domain.FeedEvent, 0, limit*2)
	eg.Go(func() error {
		// 查询发件箱
		resp, err := f.followClient.GetFollowee(ctx, &followv1.GetFolloweeRequest{Follower: uid, Limit: 10000})
		if err != nil {
			return err
		}
		followeeIDs := slice.Map(resp.FollowRelations, func(idx int, src *followv1.FollowRelation) int64 {
			return src.Followee
		})
		evts, err := f.repo.FindPullEvents(ctx, followeeIDs, timestamp, limit)
		if err != nil {
			return err
		}
		lock.Lock()
		events = append(events, evts...)
		lock.Unlock()
		return nil
	})

	eg.Go(func() error {
		evts, err := f.repo.FindPushEvents(ctx, uid, timestamp, limit)
		if err != nil {
			return err
		}
		lock.Lock()
		events = append(events, evts...)
		lock.Unlock()
		return nil
	})

	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	// 你已经查询所有的数据，现在要排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].Ctime.UnixMilli() > events[j].Ctime.UnixMilli()
	})
	return events[:slice.Min[int]([]int{int(limit), len(events)})], nil
}

func (f *feedService) GetFeedEventList(ctx context.Context, uid int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	var eg errgroup.Group
	var lock sync.Mutex
	events := make([]domain.FeedEvent, 0, limit*int64(len(f.handlerMap)))
	for _, handler := range f.handlerMap {
		h := handler
		eg.Go(func() error {
			evts, err := h.FindFeedEvents(ctx, uid, timestamp, limit)
			if err != nil {
				return err
			}
			lock.Lock()
			events = append(events, evts...)
			lock.Unlock()
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	// 你已经查询所有的数据，现在要排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].Ctime.UnixMilli() > events[j].Ctime.UnixMilli()
	})
	return events[:slice.Min[int]([]int{int(limit), len(events)})], nil
}
