package runtime

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/Luna-CY/v2ray-helper/staticfile/migrationstatic"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"io/ioutil"
	"os"
	"path/filepath"
)

var rootPath = ""

// InitRuntime 初始化运行环境
func InitRuntime(path string) error {
	if err := os.MkdirAll(filepath.Join(path, "config"), 0755); nil != err {
		return errors.New(fmt.Sprintf("初始化运行环境失败: %v", err))
	}

	mainConfigPath := filepath.Join(path, "config", "main.config.yaml")
	mainConfigExists, err := fileExists(mainConfigPath)
	if nil != err {
		return err
	}

	if !mainConfigExists {
		if err := configurator.GetDefaultMailConfig().Save(mainConfigPath); nil != err {
			return err
		}
	}

	dbPath := filepath.Join(path, "main.db")
	dbExists, err := fileExists(dbPath)
	if nil != err {
		return err
	}

	if !dbExists {
		if err := Migrate(dbPath, "1.0.0"); nil != err {
			return err
		}
	}

	rootPath = path

	return nil
}

// Migrate 数据库迁移
func Migrate(db, version string) error {
	td, err := ioutil.TempDir("", "")
	if nil != err {
		return errors.New(fmt.Sprintf("无法创建临时目录: %v", err))
	}
	defer func() {
		_ = os.RemoveAll(td)
	}()

	for _, name := range migrationstatic.AssetNames() {
		if err := os.MkdirAll(filepath.Join(td, filepath.Dir(name)), 0755); nil != err {
			return errors.New(fmt.Sprintf("无法创建migration: %v", err))
		}

		tf, err := os.OpenFile(filepath.Join(td, name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if nil != err {
			return errors.New(fmt.Sprintf("无法创建migration: %v", err))
		}

		if _, err := tf.Write(migrationstatic.MustAsset(name)); nil != err {
			return errors.New(fmt.Sprintf("无法创建migration: %v", err))
		}

		tf.Close()
	}

	sqlDb, err := sql.Open("sqlite3", db)
	if nil != err {
		return errors.New(fmt.Sprintf("无法打开数据库驱动: %v\n", err))
	}
	defer sqlDb.Close()

	driver, err := sqlite3.WithInstance(sqlDb, &sqlite3.Config{})
	if nil != err {
		return errors.New(fmt.Sprintf("无法打开Migrate驱动: %v\n", err))
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v/%v/%v", td, "migrations", version), "sqlite3", driver)
	if nil != err {
		return errors.New(fmt.Sprintf("创建Migrate失败: %v\n", err))
	}
	defer m.Close()

	if err := m.Up(); nil != err {
		if migrate.ErrNoChange != err {
			return errors.New(fmt.Sprintf("执行数据库迁移失败: %v\n", err))
		}
	}
	return nil
}

// AbsRootPath 获取项目运行绝对目录
func AbsRootPath(homeDir string) string {
	if "" != homeDir {
		if filepath.IsAbs(homeDir) {
			return homeDir
		}

		if "." != homeDir {
			absPath, err := filepath.Abs(homeDir)
			if nil == err {
				return absPath
			}
		}
	}

	// 取到执行文件所在的目录作为根目录，否则在其他目录通过文件位置运行时会找不到配置文件
	executable, err := os.Executable()
	if nil != err {
		return ""
	}

	return filepath.Dir(executable)
}

// GetRootPath 获取根目录
func GetRootPath() string {
	return rootPath
}

const CertDirName = "certs"

// GetCertificatePath 获取证书目录的路径
func GetCertificatePath() string {
	return filepath.Join(rootPath, CertDirName)
}

func fileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if nil != err {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.New(fmt.Sprintf("获取文件信息失败: %v", err))
	}

	if stat.IsDir() {
		return false, nil
	}

	return true, nil
}
