package modules

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type ColoredFormatter struct{}

func (f *ColoredFormatter) Format(entry *log.Entry) ([]byte, error) {
	// 定义颜色代码
	colorReset := "\x1b[0m"
	colorRed := "\x1b[31m"
	colorYellow := "\x1b[33m"
	colorGreen := "\x1b[32m"
	colorWhite := "\x1b[37m"

	// 获取当前时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 根据日志级别设置相应的颜色
	var levelColor string
	switch entry.Level {
	case log.InfoLevel:
		levelColor = colorGreen
	case log.WarnLevel:
		levelColor = colorYellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = colorRed
	default:
		levelColor = colorWhite
	}

	// 构建日志格式
	message := fmt.Sprintf("[%s] [%s%s%s]: %s\n", timestamp, levelColor, entry.Level, colorReset, entry.Message)
	return []byte(message), nil
}

func LogInit() {
	// 设置日志级别
	log.SetLevel(log.InfoLevel)

	// 设置日志格式为自定义格式
	log.SetFormatter(&ColoredFormatter{})

}
