package cloudreve

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestInstall(t *testing.T) {
	td := t.TempDir()

	if err := Install("linux", "amd64", "3.3.2", td, filepath.Join(td, "etc", "systemd", "system", "cloudreve.service")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "cloudreve")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "etc", "systemd", "system", "cloudreve.service")); nil != err {
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
