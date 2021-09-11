package database

import (
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var md *gorm.DB

func Init() error {
	loggerConfig := logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info, // 不输出任何日志
		Colorful:      false,
	}

	if gin.ReleaseMode == gin.Mode() {
		loggerConfig.LogLevel = logger.Silent
	}

	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), loggerConfig)

	db, err := gorm.Open(sqlite.Open(configurator.GetDbConfig().GetDatabase()), &gorm.Config{Logger: newLogger})
	if nil != err {
		return errors.New(fmt.Sprintf("打开数据库失败: %v", err))
	}

	sqlDb, err := db.DB()
	if nil != err {
		return errors.New(fmt.Sprintf("获取SQL DB失败: %v", err))
	}

	sqlDb.SetMaxOpenConns(configurator.GetDbConfig().GetMaxPoolNum())
	md = db

	return nil
}

func GetMainDb() *gorm.DB {
	if nil == md {
		panic("获取了未初始化的数据库")
	}

	return md
}
