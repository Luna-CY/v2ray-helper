package logger

import (
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

var logger *logrus.Logger

func Init(rootPath string) error {
	var dest io.Writer

	if gin.ReleaseMode == gin.Mode() {
		logPath := filepath.Join(rootPath, "logs", "main.log")
		if err := os.MkdirAll(filepath.Dir(logPath), 0755); nil != err {
			return errors.New(fmt.Sprintf("初始化日志组件失败: %v", err))
		}

		//写入文件
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return errors.New(fmt.Sprintf("初始化日志组件失败: %v", err))
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

	return nil
}

func GetLogger() *logrus.Logger {
	if nil == logger {
		panic("初始化日志失败")
	}

	return logger
}
