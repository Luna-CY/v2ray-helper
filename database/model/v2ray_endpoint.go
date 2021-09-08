package model

const (
	V2rayEndpointTransportTypeTcp = iota + 1
	V2rayEndpointTransportTypeWebSocket
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
	Level         *int    `gorm:"not null;default:0"`
	TransportType *int    `gorm:"not null;default:1"` // 传输类型：1表示TCP；2表示WebSocket
	WebSocket     *string `gorm:"not null;default:''"`
}

func (e *V2rayEndpoint) TableName() string {
	return "v2ray_endpoint"
}
