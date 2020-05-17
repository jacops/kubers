package common

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-hclog"
)

func GetLogLevel(logLevel string) (hclog.Level, error) {
	var level hclog.Level
	logLevel = strings.ToLower(strings.TrimSpace(logLevel))

	switch logLevel {
	case "trace":
		level = hclog.Trace
	case "debug":
		level = hclog.Debug
	case "notice", "info", "":
		level = hclog.Info
	case "warn", "warning":
		level = hclog.Warn
	case "err", "error":
		level = hclog.Error
	default:
		return level, fmt.Errorf("unknown log level: %s", logLevel)
	}
	return level, nil
}
