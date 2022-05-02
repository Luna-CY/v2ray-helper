package v2ray

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCheckSystem(t *testing.T) {
	if CheckSystem("windows", "386") {
		t.Fatal("测试失败")
	}

	if CheckSystem("windows", "amd64") {
		t.Fatal("测试失败")
	}

	if !CheckSystem("linux", "386") {
		t.Fatal("测试失败")
	}

	if !CheckSystem("linux", "amd64") {
		t.Fatal("测试失败")
	}

	if !CheckSystem("linux", "arm") {
		t.Fatal("测试失败")
	}

	if !CheckSystem("linux", "arm64") {
		t.Fatal("测试失败")
	}

	if CheckSystem("linux", "xxx") {
		t.Fatal("测试失败")
	}

	if CheckSystem("freebsd", "386") {
		t.Fatal("测试失败")
	}

	if CheckSystem("freebsd", "amd64") {
		t.Fatal("测试失败")
	}

	if CheckSystem("freebsd", "arm") {
		t.Fatal("测试失败")
	}

	if CheckSystem("freebsd", "arm64") {
		t.Fatal("测试失败")
	}

	if CheckSystem("darwin", "amd64") {
		t.Fatal("测试失败")
	}

	if CheckSystem("darwin", "arm64") {
		t.Fatal("测试失败")
	}

	if CheckSystem("darwin", "386") {
		t.Fatal("测试失败")
	}

	if CheckSystem("darwin", "arm") {
		t.Fatal("测试失败")
	}
}

func TestCheckExists(t *testing.T) {
	if err := CheckExists("/path/to/not-exists"); nil != err {
		t.Fatal("测试失败")
	}

	tf, err := ioutil.TempFile("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer tf.Close()

	err = CheckExists(tf.Name())
	if nil == err {
		t.Fatal("测试失败")
	}

	if !os.IsExist(err) {
		t.Fatal("测试失败")
	}

	td, err := ioutil.TempDir("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer func() {
		_ = os.RemoveAll(td)
	}()

	if err := CheckExists(td); nil != err {
		t.Fatal("测试失败")
	}
}
