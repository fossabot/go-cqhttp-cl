package main

import (
	"github.com/Redmomn/go-cqhttp-cl/modules/cmd"
	"github.com/Redmomn/go-cqhttp-cl/modules/config"
	"github.com/Redmomn/go-cqhttp-cl/proxy"
)

func main() {
	cmd.Init()
	// 启动服务器
	proxy.InitServer(config.Conf)

	select {}
}
