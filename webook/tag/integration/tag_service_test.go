package integration

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/tag/grpc"
	"gitee.com/geekbang/basic-go/webook/tag/integration/startup"
	"gitee.com/geekbang/basic-go/webook/tag/repository"
	"gitee.com/geekbang/basic-go/webook/tag/repository/cache"
	"gitee.com/geekbang/basic-go/webook/tag/repository/dao"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type TagServiceTestSuite struct {
	suite.Suite
	svc *grpc.TagServiceServer
	db  *gorm.DB
	rdb redis.Cmdable
}

func (s *TagServiceTestSuite) SetupSuite() {
	s.svc = startup.InitGRPCService()
	s.db = startup.InitTestDB()
	s.rdb = startup.InitRedis()
}

func (s *TagServiceTestSuite) TearDownSuite() {
	err := s.db.Exec("TRUNCATE TABLE `tag_bizs`").Error
	require.NoError(s.T(), err)
	// 在有外键约束的情况下，不能用 TRUNCATE
	err = s.db.Exec("DELETE FROM `tags`").Error
	require.NoError(s.T(), err)
}

func TestTagService(t *testing.T) {
	suite.Run(t, new(TagServiceTestSuite))
}

func (s *TagServiceTestSuite) TestPreload() {
	data := make([]dao.Tag, 0, 200)
	for i := 0; i < 200; i++ {
		data = append(data, dao.Tag{
			Id:   int64(i + 1),
			Name: fmt.Sprintf("tag_%d", i),
			Uid:  int64(i+1) % 3,
		})
	}
	err := s.db.Create(&data).Error
	require.NoError(s.T(), err)
	d := dao.NewGORMTagDAO(s.db)
	c := cache.NewRedisTagCache(s.rdb)
	l := startup.InitLog()
	repo := repository.NewTagRepository(d, c, l)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	err = repo.PreloadUserTags(ctx)
	require.NoError(s.T(), err)
}
