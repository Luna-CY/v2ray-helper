package util

import (
	"testing"
)

func TestGetPublicIpv4(t *testing.T) {
	_, err := GetPublicIpv4()
	if nil != err {
		t.Fatal(err)
	}
}

func TestCheckLocalPortIsAllow(t *testing.T) {
	isAllow, err := CheckLocalPortIsAllow(4123)
	if nil != err {
		t.Fatal(err)
	}

	if !isAllow {
		t.Fatal("测试失败")
	}

	isAllow, err = CheckLocalPortIsAllow(80)
	if nil != err {
		t.Fatal(err)
	}

	if isAllow {
		t.Fatal("测试失败")
	}
}
