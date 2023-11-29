package ws

import (
	"github.com/Redmomn/go-cqhttp-cl/modules/common"
	"github.com/Redmomn/go-cqhttp-cl/modules/log_imitate"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var messageChan = make(chan []byte, 200)

func connectRemoteWS(wsUrl string) {
	go func() {
		for {
			err := connect(wsUrl)
			if err != nil {
				log.Warn("连接远程websocket失败:", err)
				log.Info("将在3秒后重连")
				time.Sleep(3)
			}
		}
	}()
}

// ws
func connect(wsUrl string) error {
	// WebSocket 连接配置
	var dialer = websocket.Dialer{
		Proxy: nil,
	}

	headers := http.Header{}
	if common.Conf.AccessToken != "" {
		headers.Set("Authorization", "Bearer "+common.Conf.AccessToken)
	}

	conn, _, err := dialer.Dial(wsUrl, headers)
	if err != nil {
		return err
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	log.Info("连接到ws服务器%v", wsUrl)

	// 消息处理
	err2 := handleMessage(conn)

	return err2
}

func handleMessage(conn *websocket.Conn) error {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Warn("读取消息发生错误：", err)
			return err
		}
		messageChan <- data
		log_imitate.MsgLog(data)
	}
}
