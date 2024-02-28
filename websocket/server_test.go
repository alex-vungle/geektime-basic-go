package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	upgrader := websocket.Upgrader{}
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// 没有额外的 header 的设置
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			writer.Write([]byte("初始化 websocket 失败"))
			return
		}

		ws := &Ws{conn: conn}
		go func() {
			// 不断读取数据
			ws.ReadCycle()
		}()
		go func() {
			// 模拟输出数据
			ticker := time.NewTicker(time.Second * 3)
			for now := range ticker.C {
				// 写回时间戳
				ws.Write(now.String())
			}
		}()
	})
}

// Websocket 代表和客户端的 websocket 连接
type Ws struct {
	conn *websocket.Conn
}

func (w *Ws) ReadCycle() {
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("接收 websocket 数据失败", err)
			continue
		}
		log.Println("收到数据", string(message))
	}
}

func (w *Ws) Write(data string) {
	err := w.conn.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil {
		panic(err)
	}
}
