package log

import (
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/conf"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var mutex sync.Mutex
var instance Logger

// NewLogger 工厂方法
func NewLogger(configure conf.Configure) (Logger, error) {
	logger := &impl{configure: configure}

	return logger, logger.init()
}

// GetOrNewLogger 工厂单例
func GetOrNewLogger(configure conf.Configure) (Logger, error) {
	if nil != instance {
		return instance, nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	if nil != instance {
		return instance, nil
	}

	logger, err := NewLogger(configure)
	if nil != err {
		return nil, err
	}

	instance = logger

	return instance, nil
}

type impl struct {
	Logger

	logger    *logrus.Logger
	configure conf.Configure
}

func (s *impl) init() error {
	output := os.Stdout
	if s.configure.IsProduction() {
		path := s.configure.GetString(conf.NameLogPath)
		if !filepath.IsAbs(path) {
			path = filepath.Join(s.configure.GetHome(), path)
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); nil != err {
			return errors.New(fmt.Sprintf("无法创建日志目录: %v", err))
		}

		out, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if nil != err {
			return errors.New(fmt.Sprintf("无法打开日志文件: %v", err))
		}

		output = out
	}

	s.logger = logrus.New()
	s.logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC3339})
	s.logger.SetLevel(s.getLogLevel())
	s.logger.SetOutput(output)

	return nil
}

func (s *impl) Debug(args ...interface{}) {
	s.logger.Debug(args...)
}

func (s *impl) Debugf(format string, args ...interface{}) {
	s.logger.Debugf(format, args...)
}

func (s *impl) Info(args ...interface{}) {
	s.logger.Info(args...)
}

func (s *impl) Infof(format string, args ...interface{}) {
	s.logger.Infof(format, args...)
}

func (s *impl) Warning(args ...interface{}) {
	s.logger.Warning(args...)
}

func (s *impl) Warningf(format string, args ...interface{}) {
	s.logger.Warningf(format, args...)
}

func (s *impl) Error(args ...interface{}) {
	s.logger.Error(args...)
}

func (s *impl) Errorf(format string, args ...interface{}) {
	s.logger.Errorf(format, args...)
}

func (s *impl) getLogLevel() logrus.Level {
	switch s.configure.GetString(conf.NameLogLevel) {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warning":
		return logrus.WarnLevel
	default:
		return logrus.ErrorLevel
	}
}
