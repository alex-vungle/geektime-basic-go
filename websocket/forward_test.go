package websocket

import (
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
)

type Hub struct {
	conns *syncx.Map[string, *websocket.Conn]
}

func (h *Hub) AddConn(name string, conn *websocket.Conn) {
	h.conns.Store(name, conn)
	go func() {
		for {
			// 接收数据
			typ, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			log.Println("收到消息", string(msg))
			// 转发数据
			// 你的返回值决定了要不要继续遍历
			h.conns.Range(func(key string, value *websocket.Conn) bool {
				if key == name {
					// 我自己就不需要转发了
					return true
				}
				log.Println("转发给", key, string(msg))
				err := value.WriteMessage(typ, msg)
				if err != nil {
					// 记录日志
				}
				return true
			})
		}
	}()
}

func TestHub(t *testing.T) {
	upgrader := websocket.Upgrader{}
	hub := &Hub{conns: &syncx.Map[string, *websocket.Conn]{}}
	// 我们假定，websocket 请求发到这里
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// responseHeader 可以不传
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			writer.Write([]byte("初始化 websocket 失败"))
			return
		}
		name := request.URL.Query().Get("name")
		hub.AddConn(name, conn)
	})

	http.ListenAndServe(":8081", nil)
}
