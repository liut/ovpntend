package status

import (
	"fhyx.tech/gopak/lib/zlog"
)

func logger() zlog.Logger {
	return zlog.Get()
}
