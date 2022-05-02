package caddy

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInstallCaddy(t *testing.T) {
	td := t.TempDir()

	if err := Install("linux", "amd64", "2.4.5", td, filepath.Join(td, "etc", "caddy", "Caddyfile"), filepath.Join(td, "etc", "systemd", "system", "caddy.service")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "caddy")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "etc", "caddy", "Caddyfile")); nil != err {
		t.Fatal(err)
	}

	file, err := os.Open(filepath.Join(td, "etc", "caddy", "Caddyfile"))
	if nil != err {
		t.Fatal(err)
	}
	defer file.Close()

	result, err := ioutil.ReadAll(file)
	if nil != err {
		t.Fatal(err)
	}

	config := "import sites-enabled/*"
	if config != string(result) {
		t.Fatal(string(result))
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
