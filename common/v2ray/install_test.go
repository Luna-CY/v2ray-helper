package v2ray

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInstall(t *testing.T) {
	td, err := ioutil.TempDir("", "")
	if nil != err {
		t.Fatal(err)
	}

	defer func() {
		_ = os.RemoveAll(td)
	}()

	if err := Install("linux", "amd64", "4.41.1", td, filepath.Join(td, "share"), filepath.Join(td, "etc")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "v2ray")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "v2ctl")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "share/v2ray/geoip.dat")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "share/v2ray/geosite.dat")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "etc/systemd/system/v2ray.service")); nil != err {
		t.Fatal(err)
	}

	if err := checkFileIsExists(filepath.Join(td, "etc/systemd/system/v2ray@.service")); nil != err {
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
