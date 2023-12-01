package http

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Redmomn/go-cqhttp-cl/modules/config"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var proxyToURL string

func StartReverseProxy(listen string, ServerURL url.URL) {
	targetURL, err := url.Parse(ServerURL.String())
	if err != nil {
		log.Warn("被反向代理的URL格式非法")
		return
	}
	proxyToURL = ServerURL.String()

	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(proxyHandler)

	// 启动代理服务器
	go func() {
		err = http.ListenAndServe(listen, router)
		if err != nil {
			log.Warn("http代理服务器启动失败:", err)
		}
	}()
	log.Infof("http代理服务器已启动: http://%v -> %v", listen, targetURL)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// 获取GET请求中的URL参数
	queryParams := r.URL.Query()

	// 将URL参数转为JSON
	jsonData := make(map[string]interface{})
	for key, values := range queryParams {
		if len(values) > 0 {
			if key == "message" {
				jsonData[key] = cqFileTobase64(values[0])
			} else {
				jsonData[key] = values[0]
			}
		} else if len(values) == 0 {
			jsonData[key] = ""
		}
	}

	// 将JSON转为字节切片
	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// 构建POST请求
	req, err := http.NewRequest("POST", proxyToURL+r.URL.Path, strings.NewReader(string(jsonBytes)))
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// 设置请求头，指定内容类型为application/json
	req.Header.Set("Content-Type", "application/json")
	if config.Conf.Token != "" {
		req.Header.Set("Authorization", "Bearer "+config.Conf.Token)
	}

	// 发起HTTP请求
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 复制目标服务器的响应头到原始响应
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 将目标服务器的响应状态码写入原始响应
	w.WriteHeader(resp.StatusCode)

	// 将目标服务器的响应体写入原始响应
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func cqFileTobase64(input string) string {
	// 定义正则表达式
	re := regexp.MustCompile(`(file:///[^,\]]*)`)

	// 使用正则表达式查找所有匹配项
	matches := re.FindAllStringSubmatch(input, -1)

	// 遍历匹配项
	for _, match := range matches {
		// 获取文件路径
		filePath := match[0]

		// 读取文件内容
		fileContent, err := readFileContent(filePath)
		if err != nil {
			log.Warn("读取文件失败：", err)
			continue
		}

		// 对文件内容进行base64编码
		replacement := "base64://" + base64.StdEncoding.EncodeToString(fileContent)

		// 替换原始字符串中的匹配项
		input = strings.Replace(input, match[1], replacement, 1)
	}
	return input
}

func readFileContent(filePath string) ([]byte, error) {
	// 假设filePath是file:///开头的文件链接
	// 这里简化处理，去掉file:///前缀
	filePath = strings.TrimPrefix(filePath, "file:///")

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return content, nil
}
