package cmd

import (
	"github.com/Redmomn/go-cqhttp-cl/modules"
	"github.com/Redmomn/go-cqhttp-cl/modules/base"
	"github.com/Redmomn/go-cqhttp-cl/modules/common"
	"github.com/Redmomn/go-cqhttp-cl/modules/config"
	"github.com/Redmomn/go-cqhttp-cl/modules/log_imitate"
	"github.com/Redmomn/go-cqhttp-cl/modules/terminal"
	"github.com/Redmomn/go-cqhttp-cl/server/http"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func Init() {
	modules.LogInit()
	InitBase()

	// 加载配置文件
	log.Debug("加载配置文件")
	conf := config.Parse(base.LittleC)
	log.Infof("当前版本：%v", base.Version)
	time.Sleep(3 * time.Second)
	log_imitate.LogImit()
	common.Conf.InitConf(conf)
	// 启动服务器
	http.StartReverseProxy(conf.Local.HTTP.Address, conf.Local.HTTP.Port, conf.Remote.HTTP.URL)
}

func InitBase() {
	// 解析命令行参数
	log.Debug("解析命令行参数")
	base.Parse()
	if terminal.RunningByDoubleClick() {
		err := terminal.NoMoreDoubleClick()
		if err != nil {
			log.Errorf("遇到错误: %v", err)
			time.Sleep(time.Second * 5)
		}
		os.Exit(0)
	}
	switch {
	case base.LittleH:
		base.Help()
	}
	if base.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("进入debug模式")
	}
}
