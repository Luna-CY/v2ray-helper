package model

const (
	V2rayEndpointTransportTypeTcp = iota + 1
	V2rayEndpointTransportTypeWebSocket
	V2rayEndpointTransportTypeKcp
	V2rayEndpointTransportTypeHttp2
)

type V2rayEndpoint struct {
	Base

	Cloud    *int    `gorm:"not null"`            // 节点云服务商：1表示Vultr；2表示阿里云；3表示腾讯云；4表示华为云
	Endpoint *int    `gorm:"not null"`            // 节点位置：1表示日本；2表示香港
	Rate     *string `gorm:"not null;default:''"` // 带宽
	Remark   *string `gorm:"not null;default:''"` // 备注信息

	Host          *string `gorm:"not null"`
	Port          *int    `gorm:"not null;default:443"`
	UserId        *string `gorm:"not null"`
	AlterId       *int    `gorm:"not null;default:64"`
	UseTls        *int    `gorm:"not null;default:1"`
	Sni           *string `gorm:"not null;default:''"`
	TransportType *int    `gorm:"not null;default:1"` // 传输类型：1表示TCP；2表示WebSocket；3表示KCP；4表示HTTP2
	Tcp           *string `gorm:"not null;default:''"`
	WebSocket     *string `gorm:"not null;default:''"`
	Kcp           *string `gorm:"not null;default:''"`
	Http2         *string `gorm:"not null;default:''"`
	Grpc          *string `gorm:"not null;default:''"`
}

func (e *V2rayEndpoint) TableName() string {
	return "v2ray_endpoint"
}
