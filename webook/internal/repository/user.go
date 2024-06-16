package repository

import (
	"context"
	"database/sql"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"log"
	"time"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateEmail
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateNonZeroFields(ctx context.Context, user domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, uid int64) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

type DBConfig struct {
	DSN string
}

type CacheConfig struct {
	Addr string
}

// NewUserRepositoryV2 强耦合到了 JSON
//func NewUserRepositoryV2(cfgBytes string) *CachedUserRepository {
//	var cfg DBConfig
//	err := json.Unmarshal([]byte(cfgBytes), &cfg)
//}

// NewUserRepositoryV1 强耦合（跨层的），严重缺乏扩展性
//func NewUserRepositoryV1(dbCfg DBConfig, cCfg CacheConfig) (*CachedUserRepository, error) {
//	db, err := gorm.Open(mysql.Open(dbCfg.DSN))
//	if err != nil {
//		return nil, err
//	}
//	ud := dao.NewUserDAO(db)
//	uc := cache.NewUserCache(redis.NewClient(&redis.Options{
//		Addr: cCfg.Addr,
//	}))
//	return &CachedUserRepository{
//		dao:   ud,
//		cache: uc,
//	}, nil
//}

func NewCachedUserRepository(dao dao.UserDAO,
	c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (repo *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}

func (repo *CachedUserRepository) UpdateNonZeroFields(ctx context.Context,
	user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *CachedUserRepository) FindById(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	// 只要 err 为 nil，就返回
	if err == nil {
		return du, nil
	}

	// err 不为 nil，就要查询数据库
	// err 有两种可能
	// 1. key 不存在，说明 redis 是正常的
	// 2. 访问 redis 有问题。可能是网络有问题，也可能是 redis 本身就崩溃了

	u, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	du = repo.toDomain(u)
	//go func() {
	//	err = repo.cache.Set(ctx, du)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()

	err = repo.cache.Set(ctx, du)
	if err != nil {
		// 网络崩了，也可能是 redis 崩了
		log.Println(err)
	}
	return du, nil
}

func (repo *CachedUserRepository) FindByIdV1(ctx context.Context, uid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, uid)
	// 只要 err 为 nil，就返回
	switch err {
	case nil:
		return du, nil
	case cache.ErrKeyNotExist:
		u, err := repo.dao.FindById(ctx, uid)
		if err != nil {
			return domain.User{}, err
		}
		du = repo.toDomain(u)
		//go func() {
		//	err = repo.cache.Set(ctx, du)
		//	if err != nil {
		//		log.Println(err)
		//	}
		//}()

		err = repo.cache.Set(ctx, du)
		if err != nil {
			// 网络崩了，也可能是 redis 崩了
			log.Println(err)
		}
		return du, nil
	default:
		// 接近降级的写法
		return domain.User{}, err
	}

}

func (repo *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
