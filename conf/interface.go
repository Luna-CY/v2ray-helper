package conf

type Configure interface {
	// GetInt 获取整形值
	GetInt(key string) int

	// GetString 获取字符串型值
	GetString(key string) string

	// IsProduction 检查是否是生产环境
	IsProduction() bool

	// GetHome 获取程序主目录
	GetHome() string
}
