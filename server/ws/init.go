package ws

import (
	log "github.com/sirupsen/logrus"
	"net"
)

func StartWS(host string, port string, RemoteWSURL string) {
	connectRemoteWS(RemoteWSURL)
	StartWSServer(host, port)
	log.Infof("WS转发已启动：%v -> %v", RemoteWSURL, net.JoinHostPort(host, port))
	Clients.forwardMsg()
}
