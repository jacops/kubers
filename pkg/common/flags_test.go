package common

import (
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestCommandLogLevel(t *testing.T) {
	tests := []struct {
		level         string
		expectedLevel hclog.Level
		expectedErr   bool
	}{
		// info
		{level: "info", expectedLevel: hclog.Info, expectedErr: false},
		{level: "INFO", expectedLevel: hclog.Info, expectedErr: false},
		{level: "inFO", expectedLevel: hclog.Info, expectedErr: false},
		{level: " info ", expectedLevel: hclog.Info, expectedErr: false},
		{level: " INFO ", expectedLevel: hclog.Info, expectedErr: false},
		{level: "ofni", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "inf", expectedLevel: hclog.NoLevel, expectedErr: true},
		// notice (info)
		{level: "notice", expectedLevel: hclog.Info, expectedErr: false},
		{level: "NOTICE", expectedLevel: hclog.Info, expectedErr: false},
		{level: "nOtIcE", expectedLevel: hclog.Info, expectedErr: false},
		{level: " notice ", expectedLevel: hclog.Info, expectedErr: false},
		{level: " NOTICE ", expectedLevel: hclog.Info, expectedErr: false},
		{level: "notify", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "eciton", expectedLevel: hclog.NoLevel, expectedErr: true},
		// trace
		{level: "trace", expectedLevel: hclog.Trace, expectedErr: false},
		{level: "TRACE", expectedLevel: hclog.Trace, expectedErr: false},
		{level: "tRaCe", expectedLevel: hclog.Trace, expectedErr: false},
		{level: " trace ", expectedLevel: hclog.Trace, expectedErr: false},
		{level: " TRACE ", expectedLevel: hclog.Trace, expectedErr: false},
		{level: "tracing", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "ecart", expectedLevel: hclog.NoLevel, expectedErr: true},
		// debug
		{level: "debug", expectedLevel: hclog.Debug, expectedErr: false},
		{level: "DEBUG", expectedLevel: hclog.Debug, expectedErr: false},
		{level: "dEbUg", expectedLevel: hclog.Debug, expectedErr: false},
		{level: " debug ", expectedLevel: hclog.Debug, expectedErr: false},
		{level: " DEBUG ", expectedLevel: hclog.Debug, expectedErr: false},
		{level: "debugging", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "gubed", expectedLevel: hclog.NoLevel, expectedErr: true},
		// warn
		{level: "warn", expectedLevel: hclog.Warn, expectedErr: false},
		{level: "WARN", expectedLevel: hclog.Warn, expectedErr: false},
		{level: "wArN", expectedLevel: hclog.Warn, expectedErr: false},
		{level: " warn ", expectedLevel: hclog.Warn, expectedErr: false},
		{level: " WARN ", expectedLevel: hclog.Warn, expectedErr: false},
		{level: "warnn", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "nraw", expectedLevel: hclog.NoLevel, expectedErr: true},
		// warning (warn)
		{level: "warning", expectedLevel: hclog.Warn, expectedErr: false},
		{level: "WARNING", expectedLevel: hclog.Warn, expectedErr: false},
		{level: "wArNiNg", expectedLevel: hclog.Warn, expectedErr: false},
		{level: " warning ", expectedLevel: hclog.Warn, expectedErr: false},
		{level: " WARNING ", expectedLevel: hclog.Warn, expectedErr: false},
		{level: "warnning", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "gninnraw", expectedLevel: hclog.NoLevel, expectedErr: true},
		// err
		{level: "err", expectedLevel: hclog.Error, expectedErr: false},
		{level: "ERR", expectedLevel: hclog.Error, expectedErr: false},
		{level: "eRr", expectedLevel: hclog.Error, expectedErr: false},
		{level: " err ", expectedLevel: hclog.Error, expectedErr: false},
		{level: " ERR ", expectedLevel: hclog.Error, expectedErr: false},
		{level: "errors", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "rre", expectedLevel: hclog.NoLevel, expectedErr: true},
		// error (err)
		{level: "error", expectedLevel: hclog.Error, expectedErr: false},
		{level: "ERROR", expectedLevel: hclog.Error, expectedErr: false},
		{level: "eRrOr", expectedLevel: hclog.Error, expectedErr: false},
		{level: " error ", expectedLevel: hclog.Error, expectedErr: false},
		{level: " ERROR ", expectedLevel: hclog.Error, expectedErr: false},
		{level: "errors", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "rorre", expectedLevel: hclog.NoLevel, expectedErr: true},
		// junk
		{level: "foobar", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "junk", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "infotracedebug", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "0", expectedLevel: hclog.NoLevel, expectedErr: true},
		{level: "1", expectedLevel: hclog.NoLevel, expectedErr: true},
		// default
		{level: "", expectedLevel: hclog.Info, expectedErr: false},
		{level: " ", expectedLevel: hclog.Info, expectedErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			level, err := GetLogLevel(tt.level)
			if err != nil && !tt.expectedErr {
				t.Errorf("got error parsing log level, shouldn't have: %s", err)
			}

			if level != tt.expectedLevel {
				t.Errorf("wrong log level: got %d, expected %d", level, tt.expectedLevel)
			}
		})
	}
}
