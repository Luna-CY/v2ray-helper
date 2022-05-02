package cloudreve

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

const DefaultDbPath = "/usr/local/cloudreve/cloudreve.db"

// SetAria2 配置Aria2
func SetAria2(dbPath, rpc, token, temp string) error {
	stat, err := os.Stat(dbPath)
	if nil != err {
		return errors.New(fmt.Sprintf("获取Cloudreve数据库状态失败: %v", err))
	}

	if stat.IsDir() {
		return errors.New(fmt.Sprintf("错误的数据库路径: %v", dbPath))
	}

	if err := os.MkdirAll(temp, 0755); nil != err {
		return errors.New(fmt.Sprintf("配置Aria2失败: %v", err))
	}

	db, err := sql.Open("sqlite3", dbPath)
	if nil != err {
		return errors.New(fmt.Sprintf("无法打开数据库驱动: %v\n", err))
	}
	defer db.Close()

	settings := map[string]string{
		"aria2_token":     token,
		"aria2_rpcurl":    rpc,
		"aria2_temp_path": temp,
		"aria2_options":   "{}",
		"aria2_interval":  "60",
	}

	for option, value := range settings {
		if err := modifySetting(db, "aria2", option, value); nil != err {
			return err
		}
	}

	return nil
}

// modifySetting 更新配置表
func modifySetting(db *sql.DB, typeName, name, value string) error {
	var id int

	if err := db.QueryRow("select `id` from `settings` where `type` = ? and `name` = ?", typeName, name).Scan(&id); nil != err {
		if sql.ErrNoRows != err {
			return errors.New(fmt.Sprintf("配置Aira2失败: %v", err))
		}

		query := "insert into `settings`(`created_at`, `updated_at`, `deleted_at`, `type`, `name`, `value`) values (?, ?, null, ?, ?, ?)"
		if _, err := db.Exec(query, time.Now(), time.Now(), typeName, name, value); nil != err {
			return errors.New(fmt.Sprintf("配置Aira2失败: %v", err))
		}
	} else {
		query := "update `settings` set `value` = ?, `deleted_at` = null where `id` = ?"
		if _, err := db.Exec(query, value, id); nil != err {
			return errors.New(fmt.Sprintf("配置Aira2失败: %v", err))
		}
	}

	return nil
}
