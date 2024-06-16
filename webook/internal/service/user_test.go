package service

import (
	"context"
	"errors"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	repomocks "gitee.com/geekbang/basic-go/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_userService_Login(t *testing.T) {
	tests := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		// 预期输入
		ctx      context.Context
		email    string
		password string
		// 预期输出
		wantUser  domain.User
		wantError error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "$2a$12$Il99mx/rqx6wR6Qq5XRcAuD6o5mzE9qr80bqtlCMWesoDSwzoQWjm",
						Phone:    "1234567",
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "HelloWorld#1234",

			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "$2a$12$Il99mx/rqx6wR6Qq5XRcAuD6o5mzE9qr80bqtlCMWesoDSwzoQWjm",
				Phone:    "1234567",
			},
		},
		{
			name: "用户未找到",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "HelloWorld#1234",

			wantUser:  domain.User{},
			wantError: ErrInvalidUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("db错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "HelloWorld#1234",

			wantUser:  domain.User{},
			wantError: errors.New("db错误"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "$2a$12$Il99mx/rqx6wR6Qq5XRcAuD6o5mzE9qr80bqtlCMWesoDSwzoQWjm",
						Phone:    "1234567",
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "HelloWorld#12345",

			wantUser:  domain.User{},
			wantError: ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantError, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
