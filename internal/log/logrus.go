package log

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
)

func LogInit(fileName string) {
	//日志输出文件
	path := "./log"
	_, err := os.Stat(path) // err为nil说明目录存在
	if err != nil {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			logrus.Error("创建日志目录错误", err)
		}
	}
	// 设置在输出日志中添加文件名和方法信息
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
		// 定制文件名和函数名的输出
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := filepath.Base(frame.File)
			return frame.Function, fileName
		},
	})
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   path + fileName,
		MaxSize:    10, // megabytes
		MaxBackups: 20,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
	})
}
