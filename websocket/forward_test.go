package websocket

import (
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
)

type Hub struct {
	// syncx.Map 是我对 sync.Map 的一个简单封装
	conns *syncx.Map[string, *websocket.Conn]
}

func (h *Hub) AddConn(name string, c *websocket.Conn) {
	h.conns.Store(name, c)
	go func() {
		for {
			typ, message, err := c.ReadMessage()
			if err != nil {
				log.Println("接收 websocket 数据失败", err)
				return
			}
			// 开始转发
			h.conns.Range(func(key string, value *websocket.Conn) bool {
				if key == name {
					// 自己的消息就不要发给自己了
					return true
				}
				err1 := value.WriteMessage(typ, message)
				if err1 != nil {
					log.Println(err)
				}
				// 返回 true，确保会继续往后遍历
				return true
			})
		}
	}()
}

func TestHub(t *testing.T) {
	upgrader := websocket.Upgrader{}
	hub := &Hub{conns: &syncx.Map[string, *websocket.Conn]{}}
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// 没有额外的 header 的设置
		name := request.URL.Query().Get("name")
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			writer.Write([]byte("初始化 websocket 失败"))
			return
		}
		hub.AddConn(name, conn)
	})
	http.ListenAndServe(":8081", nil)
}
