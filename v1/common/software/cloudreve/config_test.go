package cloudreve

import (
	"database/sql"
	"path/filepath"
	"testing"
)

func TestSetAria2(t *testing.T) {
	td := t.TempDir()

	if err := SetAria2("/Users/luna/Documents/v2ray-helper-home/cloudreve.db", "http://127.0.0.1:6800", "ARIARPCACCESSTOKEN", filepath.Join(td, "cloudreve-temp")); nil != err {
		t.Fatal(err)
	}
}

func TestUpdateSetting(t *testing.T) {
	dbPath := "/Users/luna/Documents/v2ray-helper-home/cloudreve.db"

	db, err := sql.Open("sqlite3", dbPath)
	if nil != err {
		t.Fatal(err)
	}
	defer db.Close()
	defer func() {
		if _, err := db.Exec("delete from `settings` where `type` = 'test' and `name` = 'test-option'"); nil != err {
			t.Error(err)
		}
	}()

	if err := modifySetting(db, "test", "test-option", "value1"); nil != err {
		t.Fatal(err)
	}

	var value string
	if err := db.QueryRow("select `value` from `settings` where `type` = 'test' and `name` = 'test-option'").Scan(&value); nil != err {
		t.Fatal(err)
	}

	if "value1" != value {
		t.Fatal("测试失败")
	}

	if err := modifySetting(db, "test", "test-option", "value2"); nil != err {
		t.Fatal(err)
	}

	if err := db.QueryRow("select `value` from `settings` where `type` = 'test' and `name` = 'test-option'").Scan(&value); nil != err {
		t.Fatal(err)
	}

	if "value2" != value {
		t.Fatal("测试失败")
	}
}
