package web

import (
	"fhyx.tech/platform/ovpntend/pkg/zlog"
)

func logger() zlog.Logger {
	return zlog.Get()
}
