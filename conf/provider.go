package conf

import "github.com/google/wire"

var Provider = wire.NewSet(GetInstance)
