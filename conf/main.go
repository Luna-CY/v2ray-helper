package conf

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var mutex sync.Mutex
var instance Configure

// NewConfigure 工厂方法
func NewConfigure(name string, home string, paths []string) (Configure, error) {
	impl := &impl{name: name, home: home, paths: paths}

	if err := impl.init(); nil != err {
		return nil, err
	}

	return impl, nil
}

// GetOrNewConfigure 单例工厂
func GetOrNewConfigure(name string, home string, paths []string) (Configure, error) {
	if nil != instance {
		return instance, nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	if nil != instance {
		return instance, nil
	}

	configure, err := NewConfigure(name, home, paths)
	if nil != err {
		return nil, err
	}

	instance = configure

	return instance, nil
}

// GetInstance 获取配置组件单例对象
func GetInstance() Configure {
	if nil == instance {
		panic("未初始化的配置组件")
	}

	return instance
}

type impl struct {
	name  string
	home  string
	paths []string

	viper *viper.Viper
}

func (s *impl) init() error {
	s.viper = viper.New()

	s.viper.SetConfigName(s.name)
	s.viper.SetConfigType("yaml")
	for _, path := range s.paths {
		s.viper.AddConfigPath(path)
	}

	s.defaults()
	if err := s.viper.ReadInConfig(); nil != err {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return errors.New(fmt.Sprintf("无法加载配置文件: %v", err))
		}

		if err := s.viper.SafeWriteConfig(); nil != err {
			return errors.New(fmt.Sprintf("无法保存配置文件: %v", err))
		}
	}

	return nil
}

func (s *impl) GetInt(key string) int {
	return s.viper.GetInt(key)
}

func (s *impl) GetString(key string) string {
	return s.viper.GetString(key)
}

func (s *impl) IsProduction() bool {
	return "prod" == s.viper.GetString(NameEnv)
}

func (s *impl) GetHome() string {
	return s.home
}

func (s *impl) defaults() {
	s.viper.SetDefault(NameEnv, "prod")
	s.viper.SetDefault(NameListen, "127.0.0.1")
	s.viper.SetDefault(NamePort, 9000)
	s.viper.SetDefault(NameDbType, "sqlite")
	s.viper.SetDefault(NameDbDSN, "/data/home/main.db")
	s.viper.SetDefault(NameLogPath, "/var/log/home/main.log")
	s.viper.SetDefault(NameLogLevel, "error")
}
