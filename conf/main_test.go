package conf

import (
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type ImplTestSuite struct {
	suite.Suite

	td string
}

func (s *ImplTestSuite) SetupSuite() {
	s.T().Cleanup(s.Cleanup)

	td, err := os.MkdirTemp("", "")
	s.Require().Nil(err)
	s.td = td
}

func (s *ImplTestSuite) Cleanup() {
	if "" != s.td {
		_ = os.RemoveAll(s.td)
	}
}

func (s *ImplTestSuite) TestInit() {
	i := impl{name: "test.main", paths: []string{s.td}}
	err := i.init()
	s.Require().Nil(err)
	s.Require().NotNil(i.viper)
	s.Require().FileExists(filepath.Join(s.td, "test.main.yaml"))

	defaults := [][2]interface{}{{NameEnv, "prod"}, {NameListen, "127.0.0.1"}, {NamePort, 9000}, {NameDbType, "sqlite"}, {NameDbDSN, "/data/home/main.db"}, {NameLogPath, "/var/log/home/main.log"}, {NameLogLevel, "error"}}
	for _, c := range defaults {
		s.Require().Equal(i.viper.Get(c[0].(string)), c[1])
	}

	s.Require().Nil(os.Mkdir(filepath.Join(s.td, "a"), 0755))
	s.Require().Nil(os.Mkdir(filepath.Join(s.td, "b"), 0755))
	f, err := os.Create(filepath.Join(s.td, "b", "test.main.yaml"))
	s.Require().Nil(err)
	_, err = f.WriteString("env: unit\n")
	s.Require().Nil(err)
	_, err = f.WriteString("serve:\n")
	s.Require().Nil(err)
	_, err = f.WriteString("  listen: 192.168.0.1\n")
	s.Require().Nil(err)
	s.Require().Nil(f.Close())

	i = impl{name: "test.main", paths: []string{filepath.Join(s.td, "a"), filepath.Join(s.td, "b")}}
	s.Require().Nil(i.init())
	s.Require().NoFileExists(filepath.Join(s.td, "a", "test.main.yaml"))
	s.Require().Equal("unit", i.viper.GetString(NameEnv))
	s.Require().Equal("192.168.0.1", i.viper.GetString(NameListen))
}

func (s *ImplTestSuite) TestIsProduction() {
	i := impl{name: "test.main", paths: []string{s.td}}
	s.Require().Nil(i.init())

	s.Require().True(i.IsProduction())

	i.viper.Set(NameEnv, "test")
	s.Require().False(i.IsProduction())

	i.viper.Set(NameEnv, "prod")
	s.Require().True(i.IsProduction())
}

func (s *ImplTestSuite) TestNewConfigure() {
	// 工厂创建impl后调用了impl的init方法
	// 这里放到用例集来做测试，方便清理测试资源

	i, err := NewConfigure("test.main", s.td, []string{s.td})
	s.Require().Nil(err)
	s.Require().NotNil(i)
}

func (s *ImplTestSuite) TestHome() {
	i := impl{name: "test.main", home: s.td, paths: []string{s.td}}
	s.Require().Nil(i.init())

	s.Require().Equal(s.td, i.GetHome())
}

func TestImplTestSuite(t *testing.T) {
	suite.Run(t, new(ImplTestSuite))
}
