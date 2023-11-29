package terminal

import (
	"fmt"
	"github.com/Redmomn/go-cqhttp-cl/modules/base"
	"golang.org/x/sys/windows"
	"syscall"
	"time"
	"unsafe"
)

func setConsoleTitle(title string) error {
	p0, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	r1, _, err := windows.NewLazySystemDLL("kernel32.dll").NewProc("SetConsoleTitleW").Call(uintptr(unsafe.Pointer(p0)))
	if r1 == 0 {
		return err
	}
	return nil
}

// SetTitle 设置标题为 go-cqhttp `版本` `版权`
func SetTitle() {
	_ = setConsoleTitle(fmt.Sprintf("战地五暖服后端 "+base.Version+" © 2023 - %d 脑袋里进花生了", time.Now().Year()))
}
