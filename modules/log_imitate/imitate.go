package log_imitate

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stdout, "", 0)

// LogImit log模拟
func LogImit() {
	// 记录日志消息，自定义日期和时间格式
	logPrefix := "[" + time.Now().Format("2006-01-02 15:04:05") + "] [INFO]: "
	logger.SetPrefix(logPrefix)

	logger.Println("当前版本:1.14.514")
	time.Sleep(1 * time.Second)
	logger.Println("将使用 device.json 内的设备信息运行Bot.")
	time.Sleep(1 * time.Second)
	time.Sleep(1 * time.Second)
	logger.Println("Bot将在5秒后登录并开始信息处理, 按 Ctrl+C 取消.")
	time.Sleep(2 * time.Second)
	logger.Println("开始尝试登录并同步消息...")
	time.Sleep(500 * time.Millisecond)
	logger.Println("使用协议: Android Pad 8.9.63.11390")
	logger.Println("Protocol -> connect to proxy: 120.232.130.13:8080")
	time.Sleep(500 * time.Millisecond)
	logger.Println("Protocol -> device lock is disabled. HTTP API may fail.")
	time.Sleep(500 * time.Millisecond)
	logger.Println("正在检查协议更新...")
	logger.Println("登录成功 欢迎使用: robot")
	time.Sleep(500 * time.Millisecond)
	logger.Println("开始加载好友列表...")
	time.Sleep(1000 * time.Millisecond)
	logger.Println("共加载 114514 个好友.")
	logger.Println("开始加载群列表...")
	logger.Println("共加载 1919810 个群.")
	time.Sleep(500 * time.Millisecond)
	logger.Println("资源初始化完成, 开始处理信息.")
	time.Sleep(500 * time.Millisecond)
	logger.Println("アトリは、高性能ですから!")
	logger.Println("正在检查更新.")
	time.Sleep(1 * time.Millisecond)
	logger.Println("CQ WebSocket 服务器已启动: [::]:6700")
	time.Sleep(500 * time.Millisecond)
	logger.Println("CQ HTTP 服务器已启动: [::]:5700")
	time.Sleep(500 * time.Millisecond)
	logger.Println("检查更新完成. 当前已运行最新版本.")
	logger.Println("开始诊断网络情况")
	logger.Println("网络诊断完成. 未发现问题")
}

func MsgLog(data []byte) {
	var msgdata MessageData
	err := json.Unmarshal(data, &msgdata)
	if err != nil {
		logrus.Warn("消息解析出错：", err)
	}
	if msgdata.PostType == "message" {
		logger.Printf("收到群 xxx(%d) 内 %s(%d) 的消息: %s (%d)\n", msgdata.GroupID, msgdata.Sender.Nickname, msgdata.UserID, msgdata.Message, msgdata.MessageID)
	}
}

type MessageData struct {
	PostType    string `json:"post_type"`
	MessageType string `json:"message_type"`
	Time        int    `json:"time"`
	SelfID      int    `json:"self_id"`
	SubType     string `json:"sub_type"`
	RawMessage  string `json:"raw_message"`
	Sender      struct {
		Age      int    `json:"age"`
		Area     string `json:"area"`
		Card     string `json:"card"`
		Level    string `json:"level"`
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Sex      string `json:"sex"`
		Title    string `json:"title"`
		UserID   int    `json:"user_id"`
	} `json:"sender"`
	MessageID  int         `json:"message_id"`
	Anonymous  interface{} `json:"anonymous"`
	Font       int         `json:"font"`
	Message    string      `json:"message"`
	MessageSeq int         `json:"message_seq"`
	GroupID    int         `json:"group_id"`
	UserID     int         `json:"user_id"`
}
