package repository

import (
	"gitee.com/geekbang/basic-go/wirev1/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRespository(d *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: d,
	}
}
