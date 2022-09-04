package log

import (
	"os"

	"github.com/26huitailang/yogo/framework/contract"
	"github.com/26huitailang/yogo/framework/provider/log/formatter"
	"github.com/26huitailang/yogo/framework/provider/log/services"
)

// YogoTestingLog 是 Env 的具体实现
type YogoTestingLog struct {
}

// NewYogoTestingLog 测试日志，直接打印到控制台
func NewYogoTestingLog(params ...interface{}) (interface{}, error) {
	log := &services.YogoConsoleLog{}

	log.SetLevel(contract.DebugLevel)
	log.SetCtxFielder(nil)
	log.SetFormatter(formatter.TextFormatter)

	// 最重要的将内容输出到控制台
	log.SetOutput(os.Stdout)
	return log, nil
}
