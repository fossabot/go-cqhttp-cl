package main

import (
	"github.com/Redmomn/go-cqhttp-cl/modules/cmd"
	"log"
	"os"
)

// 创建自定义日志记录器
var logger = log.New(os.Stdout, "", 0)

func main() {
	cmd.Init()

}
