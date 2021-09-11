package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/Luna-CY/v2ray-helper/common/configurator"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"path"
)

type CommandArg struct {
	Up   bool
	Down bool

	Version string
}

func main() {
	ca := CommandArg{}

	flag.BoolVar(&ca.Up, "up", false, "应用迁移")
	flag.BoolVar(&ca.Down, "down", false, "撤销迁移")

	flag.StringVar(&ca.Version, "version", "", "迁移的版本号")
	flag.Parse()

	if !ca.Up && !ca.Down {
		fmt.Println("-up与-down至少指定一个动作")

		return
	}

	if ca.Up && ca.Down {
		fmt.Println("不能同时指定-up与-down参数")

		return
	}

	if "" == ca.Version {
		return
	}

	if err := configurator.Init(); nil != err {
		log.Fatalln(fmt.Sprintf("初始化配置器失败: %v", err))
	}

	if _, err := os.Stat(path.Join("migrations", ca.Version)); nil != err {
		fmt.Printf("无法获取该版本迁移数据: %v\n", err)

		return
	}

	sqlDb, err := sql.Open("sqlite3", configurator.GetDbConfig().GetDatabase())
	if nil != err {
		fmt.Printf("无法打开数据库驱动: %v\n", err)

		return
	}
	defer sqlDb.Close()

	driver, err := sqlite3.WithInstance(sqlDb, &sqlite3.Config{})
	if nil != err {
		fmt.Printf("无法打开Migrate驱动: %v\n", err)

		return
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://./migrations/%v", ca.Version), "sqlite3", driver)
	if nil != err {
		fmt.Printf("创建Migrate失败: %v\n", err)

		return
	}
	defer m.Close()

	if ca.Up {
		if err := m.Up(); nil != err {
			if migrate.ErrNoChange != err {
				fmt.Printf("执行数据库迁移失败: %v\n", err)

				return
			}
		}
	} else {
		if err := m.Down(); nil != err {
			if migrate.ErrNoChange != err {
				fmt.Printf("执行数据库迁移失败: %v\n", err)

				return
			}
		}
	}
}
