package proxy

import (
	"github.com/Redmomn/go-cqhttp-cl/modules/config"
	"github.com/Redmomn/go-cqhttp-cl/proxy/http"
	"github.com/Redmomn/go-cqhttp-cl/proxy/ws"
	"net/url"
)

func InitServer(conf config.Config) {
	if conf.HttpReverseProxy.Enable {
		for _, proxy := range conf.HttpReverseProxy.Proxys {
			proxyTo, _ := url.Parse(proxy.To)
			http.StartReverseProxy(proxy.From, *proxyTo)
		}
	}
	if conf.Websocket.Enable {
		wsUrl, _ := url.Parse(conf.Websocket.RemoteUrl)
		ws.ConnectRemoteWS(*wsUrl)
		ws.StartWSServer(conf.Websocket.Listen)

	}
}
