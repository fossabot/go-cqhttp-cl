package ws

import (
	"github.com/Redmomn/go-cqhttp-cl/modules/config"
	"github.com/Redmomn/go-cqhttp-cl/modules/log_imitate"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

var (
	onebotPostMessageChan = make(chan []byte, 200)
)

// ConnectRemoteWS 连接到远程ws，同时自动开启消息转发
func ConnectRemoteWS(wsUrl url.URL) {
	go func() {
		for {
			err := connect(wsUrl.String())
			if err != nil {
				log.Warn("连接远程websocket出错:", err)
				log.Info("将在3秒后重连")
				time.Sleep(3 * time.Second)
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

	// 工作模式为1时需要连接认证
	if config.Conf.Token != "" {
		headers.Set("Authorization", "Bearer "+config.Conf.Token)
	}

	conn, _, err := dialer.Dial(wsUrl, headers)
	if err != nil {
		return err
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	log.Info("连接到ws服务器：", wsUrl)

	// 消息处理
	err2 := handleMessage(conn)

	return err2
}

// 内部实现消息转发
func handleMessage(conn *websocket.Conn) error {
	errorChan := make(chan error)
	go func(c *chan error) {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Warn("读取消息发生错误：", err)
				errorChan <- err
				return
			}
			onebotPostMessageChan <- data
			log_imitate.MsgLog(data)
		}
	}(&errorChan)

	go func(c *chan error) {
		for {
			data := <-clientSendMessageChan

			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Warn("读取消息发生错误：", err)
				errorChan <- err
				return
			}
		}
	}(&errorChan)

	err := <-errorChan
	return err
}
