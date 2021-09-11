package v2ray

// Config 配置结构
type Config struct {
	Clients []struct {
		UserId  string `json:"user_id"`
		AlterId int    `json:"alter_id"`
	} `json:"clients"`

	V2rayPort     int `json:"v2ray_port"`
	TransportType int `json:"transport_type"`

	Tcp struct {
		Type    string `json:"type"`
		Request struct {
			Version string `json:"version"`
			Method  string `json:"method"`
			Path    string `json:"path"`
			Headers []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"headers"`
		} `json:"request"`
		Response struct {
			Version string `json:"version"`
			Status  int    `json:"status"`
			Reason  string `json:"reason"`
			Headers []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"headers"`
		} `json:"response"`
	} `json:"tcp"`
	WebSocket struct {
		Path    string `json:"path"`
		Headers []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"headers"`
	} `json:"web_socket"`
	Kcp struct {
		Type             string `json:"type"`
		Mtu              int    `json:"mtu"`
		Tti              int    `json:"tti"`
		UplinkCapacity   int    `json:"uplink_capacity"`
		DownlinkCapacity int    `json:"downlink_capacity"`
		Congestion       bool   `json:"congestion"`
		ReadBufferSize   int    `json:"read_buffer_size"`
		WriteBufferSize  int    `json:"write_buffer_size"`
	} `json:"kcp"`
	Http2 struct {
		Host string `json:"host"`
		Path string `json:"path"`
	} `json:"http2"`
}

// SetConfig 设置V2ray配置
func SetConfig(config *Config) error {
	return nil
}
