package status

import (
	"github.com/liut/ovpntend/pkg/zlog"
)

func logger() zlog.Logger {
	return zlog.Get()
}
