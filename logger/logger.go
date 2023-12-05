package logger

import (
	"github.com/beego/beego/v2/core/logs"
)

var log *logs.BeeLogger

func Init() {
	log = logs.NewLogger(100)
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/app.log","maxlines":10000000,"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)

}

func Error(message string, err error) {
	log.Error("%s -----> %v", message, err)
}

func Info(message string, data string) {
	log.Info("%s-----> %s", data)
}

func Debug(message string) {
	log.Debug(message)
}

func Alert(message string, data string) {
	log.Alert("%s -----> %v", message, data)
}

func Warning(message string, data string) {
	log.Warning("%s -----> %v", message, data)
}
