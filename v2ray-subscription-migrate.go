package main

import (
	"database/sql"
	"flag"
	"fmt"
	"gitee.com/Luna-CY/v2ray-subscription/configurator"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
	"path"
)

func main() {
	var version string

	flag.StringVar(&version, "version", "", "迁移的版本号")
	flag.Parse()

	if "" == version {
		return
	}

	if err := configurator.Init(); nil != err {
		log.Fatalln(fmt.Sprintf("初始化配置器失败: %v", err))
	}

	if _, err := os.Stat(path.Join("migrations", version)); nil != err {
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

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://./migrations/%v", version), "sqlite3", driver)
	if nil != err {
		fmt.Printf("创建Migrate失败: %v\n", err)

		return
	}
	defer m.Close()

	if err := m.Up(); nil != err {
		if migrate.ErrNoChange != err {
			fmt.Printf("执行数据库迁移失败: %v\n", err)

			return
		}
	}
}
