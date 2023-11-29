package ws

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"sync"
)

var Clients = ConnPool{}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func StartWSServer(host string, port string) {
	if !isValidIPAddress(host) {
		log.Error("启动服务器失败：非法的ip地址")
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("", handleConnections)

	go func() {
		err := http.ListenAndServe(net.JoinHostPort(host, port), router)
		if err != nil {
			log.Error("启动服务器失败:", err)
		}
	}()
}

// 解析是否合法ip地址
func isValidIPAddress(ipAddress string) bool {
	// 使用 net.ParseIP 函数来解析 IP 地址
	ip := net.ParseIP(ipAddress)

	// 如果解析成功且不为空，则认为是有效的 IP 地址
	return ip != nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("升级为ws连接时出现错误", err)
	}

	Clients.addConnection(conn)

	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

}

func (cp *ConnPool) addConnection(conn *websocket.Conn) {
	cp.dataMux.Lock()
	defer cp.dataMux.Unlock()

	cp.WSConn = append(cp.WSConn, conn)
	log.Info("接受新的websocket连接：", conn)
}

func (cp *ConnPool) removeConnection(conn *websocket.Conn) {
	cp.dataMux.RUnlock()
	cp.dataMux.Lock()
	defer cp.dataMux.Unlock()
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	for i, c := range cp.WSConn {
		if c == conn {
			cp.WSConn = append(cp.WSConn[:i], cp.WSConn[i+1:]...)
			log.Info("websocket连接断开：", conn)
			break
		}
	}
}

// TODO 需要完善检查ws状态的函数
func (cp *ConnPool) checkConnections() {
	cp.dataMux.RLock()
	defer cp.dataMux.RUnlock()

	for _, conn := range cp.WSConn {
		cp.dataMux.RUnlock()

		// TODO 这里应该有一个逻辑判断ws连接是否正常

		cp.dataMux.RLock()

		fmt.Printf("Connection %p is not connected. Removing...\n", conn)
		cp.removeConnection(conn)

	}
}

func (cp *ConnPool) closeConnection(conn *websocket.Conn) {
	cp.dataMux.Lock()
	defer cp.dataMux.Unlock()

	// Close the WebSocket connection
	if conn != nil {
		_ = conn.Close()
	}
}

// 消息转发
func (cp *ConnPool) forwardMsg() {
	var wg sync.WaitGroup
	msgData := <-messageChan
	cp.dataMux.RLock()
	defer cp.dataMux.RUnlock()

	wg.Add(1)

	go func() {
		defer wg.Done()
		for _, c := range cp.WSConn {
			err := c.WriteMessage(websocket.TextMessage, msgData)
			if err != nil {
				cp.removeConnection(c)
			}
		}
	}()

	wg.Wait()
}

type ConnPool struct {
	WSConn  []*websocket.Conn
	dataMux sync.RWMutex
}
