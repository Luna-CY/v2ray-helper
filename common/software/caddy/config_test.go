package caddy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSetConfig(t *testing.T) {
	td := t.TempDir()

	if err := SetConfig(filepath.Join(td, "sites-enabled", "localhost"), "localhost", 80, 3000, "/", false, false, false); nil != err {
		t.Fatal(err)
	}

	file, err := os.Open(filepath.Join(td, "sites-enabled", "localhost"))
	if nil != err {
		t.Fatal(err)
	}
	defer file.Close()

	result, err := ioutil.ReadAll(file)
	if nil != err {
		t.Fatal(err)
	}

	config := "localhost:80 {\n    reverse_proxy / 127.0.0.1:3000\n}"
	if config != string(result) {
		t.Fatal(string(result))
	}
}

func TestSetConfig2(t *testing.T) {
	td := t.TempDir()
	if err := SetConfig(filepath.Join(td, "sites-enabled", "localhost"), "localhost", 80, 3000, "/v2ray-path", false, true, false); nil != err {
		t.Fatal(err)
	}

	file, err := os.Open(filepath.Join(td, "sites-enabled", "localhost"))
	if nil != err {
		t.Fatal(err)
	}
	defer file.Close()

	result, err := ioutil.ReadAll(file)
	if nil != err {
		t.Fatal(err)
	}

	config := "localhost:80 {\n    reverse_proxy /v2ray-path 127.0.0.1:3000\n    reverse_proxy 127.0.0.1:5212\n}"
	if config != string(result) {
		t.Fatal(string(result))
	}
}

func TestSetConfig3(t *testing.T) {
	td := t.TempDir()
	if err := SetConfig(filepath.Join(td, "sites-enabled", "localhost"), "localhost", 443, 3000, "/", true, false, false); nil != err {
		t.Fatal(err)
	}

	file, err := os.Open(filepath.Join(td, "sites-enabled", "localhost"))
	if nil != err {
		t.Fatal(err)
	}
	defer file.Close()

	result, err := ioutil.ReadAll(file)
	if nil != err {
		t.Fatal(err)
	}

	config := "localhost:443 {\n    tls certs/localhost/cert.pem certs/localhost/private.key\n    reverse_proxy / 127.0.0.1:3000\n}"
	if config != string(result) {
		t.Fatal(string(result))
	}
}

func TestSetConfig4(t *testing.T) {
	td := t.TempDir()
	if err := SetConfig(filepath.Join(td, "sites-enabled", "localhost"), "localhost", 443, 3000, "/", true, false, true); nil != err {
		t.Fatal(err)
	}

	file, err := os.Open(filepath.Join(td, "sites-enabled", "localhost"))
	if nil != err {
		t.Fatal(err)
	}
	defer file.Close()

	result, err := ioutil.ReadAll(file)
	if nil != err {
		t.Fatal(err)
	}

	config := "localhost:443 {\n    tls certs/localhost/cert.pem certs/localhost/private.key\n    reverse_proxy / 127.0.0.1:3000 {\n        transport http {\n            versions h2c\n        }\n    }\n}"
	if config != string(result) {
		t.Fatal(string(result))
	}
}

func TestSetConfig5(t *testing.T) {
	td := t.TempDir()
	if err := SetConfig(filepath.Join(td, "sites-enabled", "localhost"), "localhost", 1234, 3000, "/", true, false, true); nil != err {
		t.Fatal(err)
	}

	file, err := os.Open(filepath.Join(td, "sites-enabled", "localhost"))
	if nil != err {
		t.Fatal(err)
	}
	defer file.Close()

	result, err := ioutil.ReadAll(file)
	if nil != err {
		t.Fatal(err)
	}

	config := "localhost:1234 {\n    tls certs/localhost/cert.pem certs/localhost/private.key\n    reverse_proxy / 127.0.0.1:3000 {\n        transport http {\n            versions h2c\n        }\n    }\n}"
	if config != string(result) {
		t.Fatal(string(result))
	}
}
