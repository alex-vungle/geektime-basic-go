//go:build wireinject

package wirev1

import (
	"gitee.com/geekbang/basic-go/wirev1/repository"
	"gitee.com/geekbang/basic-go/wirev1/repository/dao"
	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRespository, dao.NewUserDAO, InitDB)
	return &repository.UserRepository{}
}
