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

func Parse(path string) *Config {
	file, err := os.ReadFile(path)
	config := &Config{}
	if err == nil {
		err = yaml.NewDecoder(strings.NewReader(expand(string(file), os.Getenv))).Decode(config)
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
	Local  ServerConfig `yaml:"local"`
	Remote RemoteConfig `yaml:"remote"`
}

// ServerConfig 服务端配置
type ServerConfig struct {
	ImageSendDelay int `yaml:"image_send_delay"` // http和ws服务的监听地址
	WS             struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	}
	HTTP struct {
		Address string `yaml:"address"`
		Port    string `yaml:"port"`
	}
}
type RemoteConfig struct {
	AccessToken string `yaml:"access-token"`

	WS struct {
		URL string `yaml:"url"`
	}
	HTTP struct {
		URL string `yaml:"url"`
	}
}
