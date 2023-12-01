package http

import (
	"github.com/Redmomn/go-cqhttp-cl/modules/config"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func StartReverseProxy(listen string, ServerURL url.URL) {
	targetURL, err := url.Parse(ServerURL.String())
	if err != nil {
		log.Warn("被反向代理的URL格式非法")
		return
	}

	// 创建一个反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 创建一个HTTP服务器
	server := &http.Server{
		Addr: listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 在这里可以添加自定义处理逻辑，如果不需要处理，直接传递给代理服务器
			log.Debug("收到http请求: %s %s\n", r.Method, r.URL)

			// local -> remote 需要加上请求头
			if config.Conf.Token != "" {
				r.Header.Add("Authorization", "Bearer "+config.Conf.Token)

				// TODO 这里应该是要把图片上传到服务器的，并且修改cq码的图片地址，目前只做到了拦截图片的发送时间
				//if strings.Index(fmt.Sprintf("%v", r.URL), "[CQ:image") != -1 {
				//	log.Debugf("图片拦截%v秒", config.Conf.Local.ImageSendDelay)
				//	time.Sleep(time.Second * time.Duration(config.Conf.Local.ImageSendDelay))
				//}
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
	log.Infof("http代理服务器已启动: http://%v -> %v", server.Addr, targetURL)
}
