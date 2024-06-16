package repository

import (
	"context"
	"database/sql"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	cachemocks "gitee.com/geekbang/basic-go/webook/internal/repository/cache/mocks"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	daomocks "gitee.com/geekbang/basic-go/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO)

		// Input
		ctx context.Context

		uid int64
		// Expected Output
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Birthday: time.UnixMilli(100),
					AboutMe:  "自我介绍",
					Phone:    "152343434",
					Ctime:    time.UnixMilli(101)}).Return(nil)
				d.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{
					Id: uid,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "123456",
					Birthday: 100,
					AboutMe:  "自我介绍",
					Phone: sql.NullString{
						String: "152343434",
						Valid:  true,
					},
					Ctime: 101,
				}, nil)
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "152343434",
				Ctime:    time.UnixMilli(101)},
			wantErr: nil,
		},
		{
			name: "未找到用户",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				uid := 123
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), int64(uid)).Return(dao.User{}, dao.ErrRecordNotFound)
				return c, d
			},
			uid:      123,
			ctx:      context.Background(),
			wantUser: domain.User{},
			wantErr:  ErrUserNotFound,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Birthday: time.UnixMilli(100),
					AboutMe:  "自我介绍",
					Phone:    "152343434",
					Ctime:    time.UnixMilli(101)}, nil)
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "152343434",
				Ctime:    time.UnixMilli(101)},
			wantErr: nil,
		},
		{
			name: "会写缓存失败",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Birthday: time.UnixMilli(100),
					AboutMe:  "自我介绍",
					Phone:    "152343434",
					Ctime:    time.UnixMilli(101)}).Return(nil)
				d.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{
					Id: uid,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "123456",
					Birthday: 100,
					AboutMe:  "自我介绍",
					Phone: sql.NullString{
						String: "152343434",
						Valid:  true,
					},
					Ctime: 101,
				}, nil)
				return c, d
			},
			uid: 123,
			ctx: context.Background(),
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "152343434",
				Ctime:    time.UnixMilli(101)},
			wantErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc, ud := tc.mock(ctrl)
			repo := NewCachedUserRepository(ud, uc)
			user, err := repo.FindById(tc.ctx, tc.uid)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
