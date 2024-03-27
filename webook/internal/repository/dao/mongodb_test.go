package dao

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

func TestMongoDBTx(t *testing.T) {
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
	require.NoError(t, err)

	err = client.Ping(ctx, readpref.Primary())
	require.NoError(t, err)

	node, err := snowflake.NewNode(1)
	require.NoError(t, err)
	dao := NewMongoDBArticleDAOV1(client, node)
	id, err := dao.SyncWithTX(ctx, Article{
		Title:   "这是标题",
		Content: "这是内容",
	})
	require.NoError(t, err)
	t.Log(id)
}
