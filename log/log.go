package log

import (
	"fmt"
	"os"
)

type Logger struct {
	file *os.File
}

var logger = new(Logger)

func NewLogger(path string) {
	logger.CreateLogFile(path)
}

func (this *Logger) CreateLogFile(path string) {
	var file *os.File
	fileName := path + "go_p_log.log"
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(fileName)
			if err != nil {
				fmt.Println("Create file failed", err)
				return
			}
			fmt.Println("Create file success")
			this.file = file
			return
		}
	}
	fmt.Println("file exist")
	file, err = os.OpenFile(fileName, os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("open file failed")
	}
	this.file = file
}

func Debug(body ...interface{}) {
	logger.doLog(fmt.Sprint(body))
}
func Info(body ...interface{}) {
	logger.doLog(fmt.Sprint(body))
}
func Warn(body ...interface{}) {
	logger.doLog(fmt.Sprint(body))
}
func Error(body ...interface{}) {
	logger.doLog(fmt.Sprint(body))
}
func Fatal(body ...interface{}) {
	logger.doLog(fmt.Sprint(body))
}

//打印日志
func (this *Logger) doLog(body string) {
	fmt.Println(body)
	//保存日志到文件
	this.saveLog(fmt.Sprint(body, "\n"))
}
func (this *Logger) saveLog(body string) {
	_, err := this.file.WriteString(body)
	if err != nil {
		fmt.Println("save log failed:", err)
		return
	}
	//fmt.Println("save log success,total num:", n)
}
