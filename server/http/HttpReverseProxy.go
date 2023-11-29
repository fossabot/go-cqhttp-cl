package http

import (
	"fmt"
	"github.com/Redmomn/go-cqhttp-cl/modules/common"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func StartReverseProxy(host string, port string, RemoteURL string) {
	targetURL, err := url.Parse(RemoteURL)
	if err != nil {
		log.Warn("远程http地址非法")
		return
	}

	// 创建一个反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 创建一个HTTP服务器
	server := &http.Server{
		Addr: net.JoinHostPort(host, port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 在这里可以添加自定义处理逻辑，如果不需要处理，直接传递给代理服务器
			log.Debug("收到http请求: %s %s\n", r.Method, r.URL)

			if common.Conf.AccessToken != "" {
				r.Header.Add("Authorization", "Bearer "+common.Conf.AccessToken)
			}

			// TODO 这里应该是要把图片上传到服务器的，并且修改cq码的图片地址
			if strings.Index(fmt.Sprintf("%v", r.URL), "[CQ:image") != -1 {
				log.Debugf("图片拦截%v秒", common.Conf.ImageSendDelay)
				time.Sleep(time.Second * time.Duration(common.Conf.ImageSendDelay))
			}

			proxy.ServeHTTP(w, r)
		}),
	}

	// 启动代理服务器
	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Warn("http代理服务器启动失败:", err)
		}
	}()
	log.Infof("http代理服务器运行在: http://%v -> %v", server.Addr, targetURL)
}
