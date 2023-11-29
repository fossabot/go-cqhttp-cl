package common

import (
	"github.com/Redmomn/go-cqhttp-cl/modules/config"
	"sync"
)

var Conf ConfS

type ConfS struct {
	AccessToken    string
	ImageSendDelay int
	mutex          sync.Mutex
}

func (a *ConfS) InitConf(conf *config.Config) {
	defer a.mutex.Unlock()
	a.mutex.Lock()
	a.AccessToken = conf.Remote.AccessToken
	a.ImageSendDelay = conf.Local.ImageSendDelay
}
