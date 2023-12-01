package base

import (
	"flag"
	"fmt"
	"os"
)

// command flags
var (
	LittleC string // config file
	LittleH bool   // Help
)

var (
	Debug bool // 是否开启 debug 模式
)

func Parse() {
	flag.StringVar(&LittleC, "c", "conf.yml", "配置文件路径")
	flag.BoolVar(&LittleH, "h", false, "this Help")
	d := flag.Bool("D", false, "是否开启debug模式")
	flag.Parse()

	if *d {
		Debug = true
	}
}

func Help() {
	fmt.Printf(`版本：%v
使用方法：
`, Version)

	flag.PrintDefaults()
	os.Exit(0)
}
