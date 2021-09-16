package cloudreve

import (
	"path/filepath"
	"testing"
)

func TestSetAria2(t *testing.T) {
	td := t.TempDir()

	if err := SetAria2("/Users/luna/Documents/v2ray-helper-home/cloudreve.db", "http://127.0.0.1:6800", "ARIARPCACCESSTOKEN", filepath.Join(td, "cloudreve-temp")); nil != err {
		t.Fatal(err)
	}
}
