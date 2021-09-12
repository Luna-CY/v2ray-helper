package v2ray

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	tcpConfigAct1 = "{\"inbounds\":[{\"listen\":\"127.0.0.1\",\"port\":3000,\"protocol\":\"vmess\",\"settings\":{\"clients\":[{\"id\":\"c904f0ce-1385-11ec-bedb-d4619d203f36\",\"alterId\":4},{\"id\":\"test-user-id\",\"alterId\":16}]},\"streamSettings\":{\"network\":\"tcp\",\"security\":\"none\",\"tcpSettings\":{\"header\":{\"type\":\"none\"}}}}],\"outbounds\":[{\"protocol\":\"freedom\",\"settings\":{}}]}"
)

func TestSetConfig(t *testing.T) {
	tcp := &Config{V2rayPort: 3000, TransportType: TransportTypeTcp}
	tcp.Clients = append(tcp.Clients, ConfigClient{
		UserId:  "c904f0ce-1385-11ec-bedb-d4619d203f36",
		AlterId: 4,
	})
	tcp.Clients = append(tcp.Clients, ConfigClient{
		UserId:  "test-user-id",
		AlterId: 16,
	})

	tcp.Tcp.Type = "none"

	tf, err := ioutil.TempFile("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer tf.Close()

	if err := SetConfig(tf.Name(), tcp); nil != err {
		t.Fatal(err)
	}

	if err := checkFileContent(tf, tcpConfigAct1); nil != err {
		t.Fatal(err)
	}
}

func checkFileContent(file *os.File, content string) error {
	if _, err := file.Seek(io.SeekStart, 0); nil != err {
		return err
	}

	fileContent, err := ioutil.ReadAll(file)
	if nil != err {
		return err
	}

	if strings.TrimSpace(string(fileContent)) != strings.TrimSpace(content) {
		return errors.New("测试失败")
	}

	return nil
}
