package ws

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var Clients = ConnPool{}

var clientSendMessageChan = make(chan []byte, 200)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWSServer(listen string) {
	router := mux.NewRouter()
	RegisterRouters(router)

	go func() {
		err := http.ListenAndServe(listen, router)
		if err != nil {
			log.Error("启动服务器失败:", err)
		}
	}()
	log.Info("ws服务器启动成功：", listen)
	forwardToServer(&Clients)
}

func RegisterRouters(router *mux.Router) {
	router.HandleFunc("/", handleConnections)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("升级为ws连接时出现错误", err)
	}

	Clients.addConnection(&Conn{
		conn: conn,
	})

	go func(c *websocket.Conn) {
		for {
			_, data, err2 := c.ReadMessage()
			if err2 != nil {
				_ = c.Close()
				return
			}
			clientSendMessageChan <- data
		}
	}(conn)

}

// 本地ws服务器收到的消息转发到远程ws
func forwardToServer(cp *ConnPool) {
	go func() {
		for {
			data := <-onebotPostMessageChan
			for _, c := range cp.WSConn {
				c.connMux.Lock()
				err := c.conn.WriteMessage(websocket.TextMessage, data)
				c.connMux.Unlock()
				if err != nil {
					Clients.removeConnection(c)
				}
			}
		}
	}()
}
