package webserver

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallCaddy(t *testing.T) {
	td, err := ioutil.TempDir("", "")
	if nil != err {
		t.Fatal(err)
	}

	defer func() {
		_ = os.RemoveAll(td)
	}()

	if err := InstallCaddy("linux", "amd64", "2.4.5", td, filepath.Join(td, "etc", "caddy", "Caddyfile"), filepath.Join(td, "etc", "systemd", "system", "caddy.service")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "caddy")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "etc", "caddy", "Caddyfile")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "etc", "systemd", "system", "caddy.service")); nil != err {
		t.Fatal(err)
	}
}

func checkFileIsExists(path string) error {
	stat, err := os.Stat(path)
	if nil != err {
		if os.IsExist(err) {
			return nil
		}

		return err
	}

	if stat.IsDir() {
		return errors.New("错误的文件类型")
	}

	return nil
}
