package dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBArticleDAO struct {
	node    *snowflake.Node
	col     *mongo.Collection
	liveCol *mongo.Collection

	// 演示 mongodb 中的事务
	client *mongo.Client
}

func (m *MongoDBArticleDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MongoDBArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	art.Id = m.node.Generate().Int64()
	_, err := m.col.InsertOne(ctx, &art)
	return art.Id, err
}

func (m *MongoDBArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	filter := bson.D{bson.E{"id", art.Id},
		bson.E{"author_id", art.AuthorId}}
	set := bson.D{bson.E{"$set", bson.M{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   now,
	}}}
	res, err := m.col.UpdateOne(ctx, filter, set)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		// 创作者不对，说明有人在瞎搞
		return errors.New("ID 不对或者创作者不对")
	}
	return nil
}

// SyncWithTX 使用 mongodb 事务的实现
func (m *MongoDBArticleDAO) SyncWithTX(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	sess, err := m.client.StartSession()
	if err != nil {
		return 0, err
	}
	defer func() {
		sess.EndSession(ctx)
	}()

	//sess.StartTransaction()
	//sess.CommitTransaction()
	//sess.AbortTransaction()

	_, err = sess.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		if id > 0 {
			err = m.UpdateById(ctx, art)
		} else {
			id, err = m.Insert(ctx, art)
		}
		if err != nil {
			return 0, err
		}
		art.Id = id
		now := time.Now().UnixMilli()
		art.Utime = now
		// liveCol 是 INSERT or Update 语义
		filter := bson.D{bson.E{"id", art.Id},
			bson.E{"author_id", art.AuthorId}}
		set := bson.D{bson.E{"$set", art},
			bson.E{"$setOnInsert",
				bson.D{bson.E{"ctime", now}}}}
		_, err = m.liveCol.UpdateOne(ctx,
			filter, set,
			options.Update().SetUpsert(true))
		return nil, err
	})
	return id, err
}

func (m *MongoDBArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	now := time.Now().UnixMilli()
	art.Utime = now
	// liveCol 是 INSERT or Update 语义
	filter := bson.D{bson.E{"id", art.Id},
		bson.E{"author_id", art.AuthorId}}
	set := bson.D{bson.E{"$set", art},
		bson.E{"$setOnInsert",
			bson.D{bson.E{"ctime", now}}}}
	_, err = m.liveCol.UpdateOne(ctx,
		filter, set,
		options.Update().SetUpsert(true))
	return id, err
}

func (m *MongoDBArticleDAO) SyncStatus(ctx context.Context, uid int64, id int64, status uint8) error {
	filter := bson.D{bson.E{Key: "id", Value: id},
		bson.E{Key: "author_id", Value: uid}}
	sets := bson.D{bson.E{Key: "$set",
		Value: bson.D{bson.E{Key: "status", Value: status}}}}
	res, err := m.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return errors.New("ID 不对或者创作者不对")
	}
	_, err = m.liveCol.UpdateOne(ctx, filter, sets)
	return err
}

var _ ArticleDAO = &MongoDBArticleDAO{}

func NewMongoDBArticleDAO(mdb *mongo.Database, node *snowflake.Node) *MongoDBArticleDAO {
	return &MongoDBArticleDAO{
		node:    node,
		liveCol: mdb.Collection("published_articles"),
		col:     mdb.Collection("articles"),
	}
}

// NewMongoDBArticleDAOV1 利用事务的
func NewMongoDBArticleDAOV1(client *mongo.Client, node *snowflake.Node) *MongoDBArticleDAO {
	mdb := client.Database("webook")
	return &MongoDBArticleDAO{
		node:    node,
		liveCol: mdb.Collection("published_articles"),
		col:     mdb.Collection("articles"),
		client:  client,
	}
}
