package configurator

import "strings"

const defaultDatabasePath = "/data/v2ray-subscription/main.db"
const defaultPoolNum = 10

type dbConfig struct {
	Database   string `yaml:"database"`
	MaxPoolNum int    `yaml:"max-pool-num"`
}

func (d *dbConfig) GetDatabase() string {
	d.Database = strings.TrimSpace(d.Database)

	if "" == d.Database {
		return defaultDatabasePath
	}

	return d.Database
}

func (d *dbConfig) GetMaxPoolNum() int {
	if 0 == d.MaxPoolNum {

		return defaultPoolNum
	}

	return d.MaxPoolNum
}
