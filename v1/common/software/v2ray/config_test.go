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
	tcpConfigAct1  = "{\"inbounds\":[{\"listen\":\"127.0.0.1\",\"port\":3000,\"protocol\":\"vmess\",\"settings\":{\"clients\":[{\"id\":\"c904f0ce-1385-11ec-bedb-d4619d203f36\",\"alterId\":4},{\"id\":\"test-user-id\",\"alterId\":16}]},\"streamSettings\":{\"network\":\"tcp\",\"security\":\"none\",\"tcpSettings\":{\"header\":{\"type\":\"none\"}}}}],\"outbounds\":[{\"protocol\":\"freedom\",\"settings\":{}}]}"
	tcpConfigAct2  = "{\"inbounds\":[{\"listen\":\"127.0.0.1\",\"port\":3000,\"protocol\":\"vmess\",\"settings\":{\"clients\":[{\"id\":\"c904f0ce-1385-11ec-bedb-d4619d203f36\",\"alterId\":4},{\"id\":\"test-user-id\",\"alterId\":16}]},\"streamSettings\":{\"network\":\"tcp\",\"security\":\"none\",\"tcpSettings\":{\"header\":{\"type\":\"http\",\"request\":{\"version\":\"1.1\",\"method\":\"POST\",\"path\":[\"/\"],\"headers\":{\"key1\":[\"val1\"],\"key2\":[\"val2\",\"val3\"]}},\"response\":{\"version\":\"1.1\",\"status\":\"204\",\"reason\":\"Success\",\"headers\":{\"key1\":[\"val1\"],\"key2\":[\"val2\",\"val3\"]}}}}}}],\"outbounds\":[{\"protocol\":\"freedom\",\"settings\":{}}]}"
	wsConfigAct1   = "{\"inbounds\":[{\"listen\":\"127.0.0.1\",\"port\":3000,\"protocol\":\"vmess\",\"settings\":{\"clients\":[{\"id\":\"c904f0ce-1385-11ec-bedb-d4619d203f36\",\"alterId\":4},{\"id\":\"test-user-id\",\"alterId\":16}]},\"streamSettings\":{\"network\":\"ws\",\"security\":\"none\",\"wsSettings\":{\"path\":\"/example-path\",\"headers\":{}}}}],\"outbounds\":[{\"protocol\":\"freedom\",\"settings\":{}}]}"
	wsConfigAct2   = "{\"inbounds\":[{\"listen\":\"127.0.0.1\",\"port\":3000,\"protocol\":\"vmess\",\"settings\":{\"clients\":[{\"id\":\"c904f0ce-1385-11ec-bedb-d4619d203f36\",\"alterId\":4},{\"id\":\"test-user-id\",\"alterId\":16}]},\"streamSettings\":{\"network\":\"ws\",\"security\":\"none\",\"wsSettings\":{\"path\":\"/example-path\",\"headers\":{\"key1\":\"val1\"}}}}],\"outbounds\":[{\"protocol\":\"freedom\",\"settings\":{}}]}"
	kcpConfigAct   = "{\"inbounds\":[{\"listen\":\"127.0.0.1\",\"port\":3000,\"protocol\":\"vmess\",\"settings\":{\"clients\":[{\"id\":\"c904f0ce-1385-11ec-bedb-d4619d203f36\",\"alterId\":4},{\"id\":\"test-user-id\",\"alterId\":16}]},\"streamSettings\":{\"network\":\"kcp\",\"security\":\"none\",\"kcpSettings\":{\"header\":{\"type\":\"none\"},\"mtu\":1350,\"tti\":50,\"uplinkCapacity\":10,\"downlinkCapacity\":10,\"congestion\":false,\"readBufferSize\":5,\"writeBufferSize\":5}}}],\"outbounds\":[{\"protocol\":\"freedom\",\"settings\":{}}]}"
	http2ConfigAct = "{\"inbounds\":[{\"listen\":\"127.0.0.1\",\"port\":3000,\"protocol\":\"vmess\",\"settings\":{\"clients\":[{\"id\":\"c904f0ce-1385-11ec-bedb-d4619d203f36\",\"alterId\":4},{\"id\":\"test-user-id\",\"alterId\":16}]},\"streamSettings\":{\"network\":\"http\",\"security\":\"none\",\"httpSettings\":{\"host\":[\"host1.com\",\"host2.com\"],\"path\":\"/\"}}}],\"outbounds\":[{\"protocol\":\"freedom\",\"settings\":{}}]}"
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

func TestSetConfig2(t *testing.T) {
	tcp := &Config{V2rayPort: 3000, TransportType: TransportTypeTcp}
	tcp.Clients = append(tcp.Clients, ConfigClient{
		UserId:  "c904f0ce-1385-11ec-bedb-d4619d203f36",
		AlterId: 4,
	})
	tcp.Clients = append(tcp.Clients, ConfigClient{
		UserId:  "test-user-id",
		AlterId: 16,
	})

	tcp.Tcp.Type = "http"
	tcp.Tcp.Request.Version = "1.1"
	tcp.Tcp.Request.Method = "POST"
	tcp.Tcp.Request.Path = "/"
	tcp.Tcp.Request.Headers = append(tcp.Tcp.Request.Headers, ConfigHeader{Key: "key1", Value: "val1"}, ConfigHeader{Key: "key2", Value: "val2;;;val3"})

	tcp.Tcp.Response.Version = "1.1"
	tcp.Tcp.Response.Status = "204"
	tcp.Tcp.Response.Reason = "Success"
	tcp.Tcp.Response.Headers = append(tcp.Tcp.Response.Headers, ConfigHeader{Key: "key1", Value: "val1"}, ConfigHeader{Key: "key2", Value: "val2;;;val3"})

	tf, err := ioutil.TempFile("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer tf.Close()

	if err := SetConfig(tf.Name(), tcp); nil != err {
		t.Fatal(err)
	}

	if err := checkFileContent(tf, tcpConfigAct2); nil != err {
		t.Fatal(err)
	}
}

func TestSetConfig3(t *testing.T) {
	ws := &Config{V2rayPort: 3000, TransportType: TransportTypeWebSocket}
	ws.Clients = append(ws.Clients, ConfigClient{
		UserId:  "c904f0ce-1385-11ec-bedb-d4619d203f36",
		AlterId: 4,
	})
	ws.Clients = append(ws.Clients, ConfigClient{
		UserId:  "test-user-id",
		AlterId: 16,
	})

	ws.WebSocket.Path = "/example-path"

	tf, err := ioutil.TempFile("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer tf.Close()

	if err := SetConfig(tf.Name(), ws); nil != err {
		t.Fatal(err)
	}

	if err := checkFileContent(tf, wsConfigAct1); nil != err {
		t.Fatal(err)
	}
}

func TestSetConfig4(t *testing.T) {
	ws := &Config{V2rayPort: 3000, TransportType: TransportTypeWebSocket}
	ws.Clients = append(ws.Clients, ConfigClient{
		UserId:  "c904f0ce-1385-11ec-bedb-d4619d203f36",
		AlterId: 4,
	})
	ws.Clients = append(ws.Clients, ConfigClient{
		UserId:  "test-user-id",
		AlterId: 16,
	})

	ws.WebSocket.Path = "/example-path"
	ws.WebSocket.Headers = append(ws.WebSocket.Headers, ConfigHeader{Key: "key1", Value: "val1"})

	tf, err := ioutil.TempFile("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer tf.Close()

	if err := SetConfig(tf.Name(), ws); nil != err {
		t.Fatal(err)
	}

	if err := checkFileContent(tf, wsConfigAct2); nil != err {
		t.Fatal(err)
	}
}

func TestSetConfig5(t *testing.T) {
	kcp := &Config{V2rayPort: 3000, TransportType: TransportTypeKcp}
	kcp.Clients = append(kcp.Clients, ConfigClient{
		UserId:  "c904f0ce-1385-11ec-bedb-d4619d203f36",
		AlterId: 4,
	})
	kcp.Clients = append(kcp.Clients, ConfigClient{
		UserId:  "test-user-id",
		AlterId: 16,
	})

	kcp.Kcp.Type = "none"
	kcp.Kcp.Mtu = 1350
	kcp.Kcp.Tti = 50
	kcp.Kcp.Congestion = false
	kcp.Kcp.UpLinkCapacity = 10
	kcp.Kcp.DownLinkCapacity = 10
	kcp.Kcp.ReadBufferSize = 5
	kcp.Kcp.WriteBufferSize = 5

	tf, err := ioutil.TempFile("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer tf.Close()

	if err := SetConfig(tf.Name(), kcp); nil != err {
		t.Fatal(err)
	}

	if err := checkFileContent(tf, kcpConfigAct); nil != err {
		t.Fatal(err)
	}
}

func TestSetConfig6(t *testing.T) {
	http2 := &Config{V2rayPort: 3000, TransportType: TransportTypeHttp2}
	http2.Clients = append(http2.Clients, ConfigClient{
		UserId:  "c904f0ce-1385-11ec-bedb-d4619d203f36",
		AlterId: 4,
	})
	http2.Clients = append(http2.Clients, ConfigClient{
		UserId:  "test-user-id",
		AlterId: 16,
	})

	http2.Http2.Host = "host1.com,host2.com"
	http2.Http2.Path = "/"

	tf, err := ioutil.TempFile("", "")
	if nil != err {
		t.Fatal(err)
	}
	defer tf.Close()

	if err := SetConfig(tf.Name(), http2); nil != err {
		t.Fatal(err)
	}

	if err := checkFileContent(tf, http2ConfigAct); nil != err {
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
