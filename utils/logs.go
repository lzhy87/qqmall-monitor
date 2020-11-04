package utils

import (
	"github.com/astaxie/beego/logs"
)

//Logs 定义Logs 全局变量
var Logs *logs.BeeLogger

func init() {
	Logs = logs.NewLogger()
	Logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/monitor.log","maxdays": 2,"separate":["error", "warning", "info"]}`)
	Logs.SetLogger(logs.AdapterConsole)
	Logs.EnableFuncCallDepth(true)
	// Logs.Async()

}
