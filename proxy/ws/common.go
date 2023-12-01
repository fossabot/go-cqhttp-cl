package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"sync"
)

func (cp *ConnPool) addConnection(conn *Conn) {
	cp.dataMux.Lock()
	defer cp.dataMux.Unlock()

	cp.WSConn = append(cp.WSConn, conn)
	log.Info("接受新的websocket连接：", conn.conn.RemoteAddr())
}

func (cp *ConnPool) removeConnection(conn *Conn) {
	cp.dataMux.Lock()
	defer cp.dataMux.Unlock()
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn.conn)

	for i, c := range cp.WSConn {
		if c.conn == conn.conn {
			cp.WSConn = append(cp.WSConn[:i], cp.WSConn[i+1:]...)
			log.Info("websocket连接断开：", conn.conn.RemoteAddr())
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

type ConnPool struct {
	WSConn  []*Conn
	dataMux sync.RWMutex
}

type Conn struct {
	conn    *websocket.Conn
	connMux sync.RWMutex
}
