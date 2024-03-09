package simpleim

import (
	"context"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"
)

// websocket 的网关
type WsGateway struct {
	svc        *Service
	client     sarama.Client
	l          logger.LoggerV1
	conns      *syncx.Map[int64, *Conn]
	instanceID string
	upgrader   *websocket.Upgrader
}

func (s *WsGateway) Start(addr string) error {
	// 启动。然后监听端口，接收websocket 请求
	// 我们假定，websocket 请求发到这里
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.wsHandler)
	err := s.subscribeMsg()
	if err != nil {
		// 记录日志
	}
	// 启动
	return http.ListenAndServe(addr, mux)
}

func (s *WsGateway) wsHandler(writer http.ResponseWriter, request *http.Request) {
	// responseHeader 可以不传
	c, err := s.upgrader.Upgrade(writer, request, nil)
	if err != nil {
		writer.Write([]byte("初始化 websocket 失败"))
		return
	}
	// Uid 怎么搞到呢?
	// 从 JWT Token 或者 session 里面搞到
	uid := s.Uid(request)
	conn := &Conn{Conn: c}
	s.conns.Store(uid, conn)
	// 在这里，我要填充具体的内容
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// 记录日志
				return
			}
			var msg Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				// 记录日志
				// 消息格式不对，但是websocket
				continue
			}
			// 理论上来说，这就要干嘛转发到后端
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				err = s.svc.Receive(ctx, uid, msg)
				cancel()
				if err != nil {
					// 后端服务处理失败了
					// 告诉前端，我失败了
					err = conn.Send(Message{Type: "result", Content: "FAILED", Seq: msg.Seq})
					if err != nil {
						// 记录日志
					}
				}
			}()
		}
	}()
}

// 启动消费者，监听 Kafka
func (s *WsGateway) subscribeMsg() error {
	cg, err := sarama.NewConsumerGroupFromClient(s.instanceID, s.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{eventName},
			saramax.NewHandler[Event](s.l, s.Consume))
		if err != nil {
			// 记录日志
			return
		}
	}()
	return err
}

func (s *WsGateway) Uid(req *http.Request) int64 {
	// 模拟从 header 拿出来
	uidStr := req.Header.Get("uid")
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	return uid
}

func (s *WsGateway) Consume(msg *sarama.ConsumerMessage, event Event) error {
	// 我要消费
	conn, ok := s.conns.Load(event.Receiver)
	if !ok {
		// 不需要转发
		return nil
	}
	return conn.Send(event.Msg)
}

type Conn struct {
	*websocket.Conn
}

func (c *Conn) Send(msg Message) error {
	val, _ := json.Marshal(msg)
	return c.Conn.WriteMessage(websocket.TextMessage, val)
}

// Message 前后端交互的数据格式
type Message struct {
	// 前端的序列号
	Seq string `json:"seq"`
	// 标记是什么类型的消息
	// 比如说图片，视频
	// {"type": "image", content:"http://myimage"}
	Type string `json:"type"`
	// 内容肯定有
	Content string `json:"content"`
	// 你发给谁？
	// channel id
	Cid int64
}
