//go:build !windows

package terminal

import (
	"fmt"
	"github.com/Redmomn/go-cqhttp-cl/modules/base"
	"time"
)

// SetTitle 设置标题为 go-cqhttp `版本` `版权`
func SetTitle() {
	fmt.Printf("\033]0;战地五暖服后端 "+base.Version+" © 2023 - %d 脑袋里进花生了"+"\007", time.Now().Year())
}
