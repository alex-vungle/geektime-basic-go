package simpleim

import (
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type WsGateway struct {
	// 连接了这个实例的客户端
	// 这里我们用 uid 作为 key
	// 实践中要考虑到不同的设备，
	// 那么这个 key 可能是一个复合结构，例如 uid + 设备
	conns *syncx.Map[int64, *Conn]
	svc   *IMService

	client     sarama.Client
	instanceId string
}

func (g *WsGateway) Start(addr string) error {

}

func (g *WsGateway) subscribeMsg() error {
	cg, err := sarama.NewConsumerGroupFromClient(g.instanceId,
		g.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{eventName},
			saramax.NewHandler[Event](logger.NewNoOpLogger(), g.consume))
		if err != nil {
			log.Println("退出监听消息循环", err)
		}
	}()
	return nil
}

func (g *WsGateway) wsHandler(writer http.ResponseWriter, request *http.Request) {

}

// Uid 一般是从 jwt token 或者 session 里面取出来
// 这里模拟从 header 里面读取出来
func (g *WsGateway) Uid(req *http.Request) int64 {
	uidStr := req.Header.Get("uid")
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	return uid
}

func (g *WsGateway) consume(msg *sarama.ConsumerMessage, evt Event) error {
	panic("implement me")
}

// Conn 稍微做一个封装
type Conn struct {
	*websocket.Conn
}

func (c *Conn) Send(msg Message) error {
	val, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.WriteMessage(websocket.TextMessage, val)
}

type Message struct {
	// 发过来的消息的序列号
	// 用于前后端关联消息
	Seq string
	// 用来标识不同的消息类型
	// 文本消息，视频消息
	// 系统消息（后端往前端发的，跟 IM 本身管理有关的消息）
	Type    string
	Content string
	// 聊天 ID，注意，正常来说这里不是记录目标用户 ID
	// 而是记录代表了这个聊天的 ID
	Cid int64
}
