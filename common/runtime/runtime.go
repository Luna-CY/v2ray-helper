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

func InitRuntime(rootPath string) error {
	if err := os.MkdirAll(filepath.Join(rootPath, "config"), 0755); nil != err {
		return errors.New(fmt.Sprintf("初始化运行环境失败: %v", err))
	}

	mainConfigPath := filepath.Join(rootPath, "config", "main.prod.config.yaml")
	mainConfigExists, err := fileExists(mainConfigPath)
	if nil != err {
		return err
	}

	if !mainConfigExists {
		mainConfigFile, err := os.OpenFile(mainConfigPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if nil != err {
			return errors.New(fmt.Sprintf("无法打开文件: %v", err))
		}
		defer mainConfigFile.Close()

		if _, err := mainConfigFile.WriteString(configurator.DefaultMainConfigContent); nil != err {
			return errors.New(fmt.Sprintf("无法写入文件: %v", err))
		}
	}

	dbPath := filepath.Join(rootPath, "main.db")
	dbExists, err := fileExists(dbPath)
	if nil != err {
		return err
	}

	if !dbExists {
		if err := Migrate(dbPath, "1.0.0"); nil != err {
			return err
		}
	}

	return nil
}

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

// GetRootPath 获取项目运行根目录
func GetRootPath(homeDir string) string {
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
