package mongodb

import (
	"context"
	"fmt"
	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/chenmingyong0423/go-mongox/builder/update"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

// 这是一个我们同学写的改进 mongoDB 操作的库，我觉得蛮好用的
type MongoXTestSuite struct {
	suite.Suite
	client *mongo.Client
	col    *mongo.Collection
}

func (s *MongoXTestSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			fmt.Println(evt.Command)
		},
	}
	opts := options.Client().
		ApplyURI("mongodb://root:example@localhost:27017/").
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	require.NoError(s.T(), err)
	err = client.Ping(ctx, readpref.Primary())
	require.NoError(s.T(), err)
	s.client = client
	s.col = s.client.Database("webook").Collection("article")
}

func (s *MongoXTestSuite) TestCRUD() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// 转换成它自己的 Collection 表达，利用泛型这样可以减少代码错误
	artCol := mongox.NewCollection[Article](s.col)
	// 有编译器保护
	id, err := artCol.Creator().InsertOne(ctx, &Article{
		Id:       123,
		Title:    "我的标题",
		Content:  "我的内容",
		AuthorId: 123,
	})
	require.NoError(s.T(), err)
	s.T().Log(id)
	// 更新操作
	ures, err := artCol.Updater().
		// Updates 是指定更新什么字段
		Updates(update.Set("title", "新的标题")).
		Filter(query.Eq("author_id", 123)).
		UpdateOne(ctx)
	require.NoError(s.T(), err)
	s.T().Log(ures.ModifiedCount)

	// 根据 ID 和 author_id 来查找
	flt := query.And(query.Eq("author_id", 123),
		query.Eq("title", "新的标题"))
	art, err := artCol.Finder().Filter(flt).FindOne(ctx)
	require.NoError(s.T(), err)
	// art 直接就是 Article 类型，非常省事
	s.T().Log(art.Title)
	dres, err := artCol.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	require.NoError(s.T(), err)
	s.T().Log(dres.DeletedCount)
}

func TestMongoX(t *testing.T) {
	suite.Run(t, new(MongoXTestSuite))
}
