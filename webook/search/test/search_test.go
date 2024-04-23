package test

import (
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/search/domain"
	"gitee.com/geekbang/basic-go/webook/search/events"
	"gitee.com/geekbang/basic-go/webook/search/repository/dao"
	"gitee.com/geekbang/basic-go/webook/search/service"
	"gitee.com/geekbang/basic-go/webook/search/test/startup"
	"github.com/IBM/sarama"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SearchTestSuite struct {
	suite.Suite
	svc        service.SearchService
	producer   sarama.SyncProducer
	articleDao dao.ArticleDAO
	anyDao     dao.AnyDAO
	likeDao    dao.LikeDAO
	collectDao dao.CollectDAO
	es         *elastic.Client
}

func (s *SearchTestSuite) SetupSuite() {
	svc, p, anyDao, articleDao, likeDao, collectDao, es := startup.InitTestSvc()
	s.svc = svc
	s.producer = p
	s.anyDao = anyDao
	s.articleDao = articleDao
	s.likeDao = likeDao
	s.collectDao = collectDao
	s.es = es
}

// 测试同步点赞和收藏数
func (s *SearchTestSuite) TestSync() {
	// 往队列发消息
	time.Sleep(20 * time.Second)
	datas := []events.InteractiveEvent{
		{
			Type:  1,
			Uid:   2,
			Biz:   "article",
			BizId: 1,
		},
		{
			Type:  1,
			Uid:   2,
			Biz:   "article",
			BizId: 2,
		},
		{
			Type:  2,
			Uid:   2,
			Biz:   "article",
			BizId: 1,
		},
		{
			Type:  3,
			Uid:   2,
			Biz:   "article",
			BizId: 1,
		},
	}

	for _, data := range datas {
		val, err := json.Marshal(data)
		require.NoError(s.T(), err)
		_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
			Topic: events.InteractiveTopic,
			Value: sarama.ByteEncoder(val),
		})
		require.NoError(s.T(), err)
	}
	// 等待消息处理完成
	// 查询es
	time.Sleep(10 * time.Second)
	likes, err := s.likeDao.Search(context.Background(), 2, "article")
	require.NoError(s.T(), err)
	collects, err := s.collectDao.Search(context.Background(), 2, "article")
	require.NoError(s.T(), err)

	assert.Equal(s.T(), 1, len(likes))
	assert.Equal(s.T(), []int64{2}, likes)
	assert.Equal(s.T(), 1, len(collects))
	assert.Equal(s.T(), []int64{1}, collects)
}

func (s *SearchTestSuite) TestSearch() {
	err := s.articleDao.InputArticle(context.Background(), dao.Article{
		Id:      1,
		Title:   "test1",
		Status:  2,
		Content: "csswww文章1",
	})
	err = s.articleDao.InputArticle(context.Background(), dao.Article{
		Id:      2,
		Title:   "test2",
		Status:  2,
		Content: "csss文章2",
	})
	err = s.articleDao.InputArticle(context.Background(), dao.Article{
		Id:      3,
		Title:   "test3",
		Status:  2,
		Content: "css文章3",
	})
	err = s.articleDao.InputArticle(context.Background(), dao.Article{
		Id:      4,
		Title:   "test4",
		Status:  2,
		Content: "csswww文章4",
	})
	err = s.articleDao.InputArticle(context.Background(), dao.Article{
		Id:      5,
		Title:   "test5",
		Status:  2,
		Content: "csswww文章5",
	})
	collect1 := events.InteractiveEvent{
		Uid:   1,
		Biz:   "article",
		BizId: 1,
	}
	data1, err := json.Marshal(collect1)
	require.NoError(s.T(), err)
	like2 := events.InteractiveEvent{
		Uid:   1,
		Biz:   "article",
		BizId: 2,
	}
	data2, err := json.Marshal(like2)
	require.NoError(s.T(), err)
	tag3 := tag{
		Tags:  "tag1",
		UID:   1,
		Biz:   "article",
		BizID: 3,
	}
	data3, err := json.Marshal(tag3)
	require.NoError(s.T(), err)
	err = s.anyDao.Input(context.Background(), "collect_index", "1_article_1", string(data1))
	require.NoError(s.T(), err)
	err = s.anyDao.Input(context.Background(), "like_index", "1_article_2", string(data2))
	require.NoError(s.T(), err)
	err = s.anyDao.Input(context.Background(), "tags_index", "1_article_3", string(data3))
	require.NoError(s.T(), err)
	time.Sleep(1 * time.Second)
	res, err := s.svc.Search(context.Background(), 1, "tag1 test4")
	require.NoError(s.T(), err)
	articles := res.Articles
	assert.Equal(s.T(), 4, len(articles))
	assert.Equal(s.T(), []domain.Article{
		{
			Id:      4,
			Title:   "test4",
			Status:  2,
			Content: "csswww文章4",
		},
		{
			Id:      1,
			Title:   "test1",
			Status:  2,
			Content: "csswww文章1",
		},
		{
			Id:      3,
			Title:   "test3",
			Status:  2,
			Content: "css文章3",
		},
		{
			Id:      2,
			Title:   "test2",
			Status:  2,
			Content: "csss文章2",
		},
	}, articles)
}

func TestSearch(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}

type tag struct {
	Tags  string `json:"tags"`
	UID   int64  `json:"uid"`
	Biz   string `json:"biz"`
	BizID int64  `json:"biz_id"`
}
