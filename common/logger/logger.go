package logger

import (
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var logger *logrus.Logger

func Init() error {
	var dest io.Writer

	if gin.ReleaseMode == gin.Mode() {
		//写入文件
		file, err := os.OpenFile(configurator.GetMainConfig().GetLogPath(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("err", err)
		}

		dest = file
	} else {
		dest = os.Stdout
	}

	//实例化
	logger = logrus.New()

	//设置输出
	logger.SetOutput(dest)

	//设置日志级别
	logger.SetLevel(configurator.GetMainConfig().GetLogLevel())

	//设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{})

	return nil
}

func GetLogger() *logrus.Logger {
	if nil == logger {
		if err := Init(); nil != err {
			panic("初始化日志失败")
		}
	}

	return logger
}
