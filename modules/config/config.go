package config

import (
	"bufio"
	_ "embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
)

//go:embed default_config.yaml
var defaultConfig string

var Conf Config

func Parse(path string) Config {
	file, err := os.ReadFile(path)
	config := Config{}
	if err == nil {
		err = yaml.NewDecoder(strings.NewReader(expand(string(file), os.Getenv))).Decode(&config)
		if err != nil {
			log.Fatal("配置文件不合法!", err)
		}
	} else {
		generateConfig()
		os.Exit(0)
	}
	return config
}

// generateConfig 生成配置文件
func generateConfig() {
	fmt.Println("未找到配置文件，正在为您生成配置文件中！")
	sb := strings.Builder{}
	sb.WriteString(defaultConfig)
	_ = os.WriteFile("conf.yml", []byte(sb.String()), 0o644)
	fmt.Println("默认配置文件已生成，请修改 conf.yml 后重新启动!")
	input := bufio.NewReader(os.Stdin)
	_, _ = input.ReadString('\n')
}

// expand 替换环境变量
func expand(s string, mapping func(string) string) string {
	r := regexp.MustCompile(`\${([a-zA-Z_]+[a-zA-Z0-9_:/.]*)}`)
	return r.ReplaceAllStringFunc(s, func(s string) string {
		s = strings.Trim(s, "${}")
		before, after, ok := strings.Cut(s, ":")
		m := mapping(before)
		if ok && m == "" {
			return after
		}
		return m
	})
}

// Config 解析配置的结构体
type Config struct {
	Token            string          `yaml:"token"`
	Websocket        WebSocketConfig `yaml:"websocket"`
	HttpReverseProxy ProxyConfig     `yaml:"http_reverse_proxy"`
	TcpProxy         ProxyConfig     `yaml:"tcp_proxy"`
}

type WebSocketConfig struct {
	Enable    bool   `yaml:"enable"`
	RemoteUrl string `yaml:"url"`
	Listen    string `yaml:"listen"`
}

type ProxyConfig struct {
	Enable bool           `yaml:"enable"`
	Proxys []ProxysConfig `yaml:"proxys"`
}

type ProxysConfig struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}
