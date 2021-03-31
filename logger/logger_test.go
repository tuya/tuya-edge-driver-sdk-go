//
// Copyright (c) 2018 Cavium
//
// SPDX-License-Identifier: Apache-2.0
//

package logger

import (
	"testing"

	"go.uber.org/zap"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"

	"github.com/stretchr/testify/assert"
)

func TestIsValidLogLevel(t *testing.T) {
	var tests = []struct {
		level string
		res   bool
	}{
		{models.TraceLog, true},
		{models.DebugLog, true},
		{models.InfoLog, true},
		{models.WarnLog, true},
		{models.ErrorLog, true},
		{"EERROR", false},
		{"ERRORR", false},
		{"INF", false},
	}
	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			r := isValidLogLevel(tt.level)
			if r != tt.res {
				t.Errorf("Level %s labeled as %v and should be %v",
					tt.level, r, tt.res)
			}
		})
	}
}

type Demo struct {
	User string `json:"user"`
}

func TestLogLevel(t *testing.T) {
	expectedLogLevel := models.DebugLog
	lc := NewClient("testService", expectedLogLevel)
	lc.Debug("message", Demo{User: "123"})
	lc.Info("message", []string{"1", "2"})
	lc.Warn("message", map[string]interface{}{"1": 233, "2": "23333"})
	lc.Error("message")
	assert.Equal(t, expectedLogLevel, lc.LogLevel())
}

func TestNewZapLogger(t *testing.T) {
	s := []string{
		"hello debug",
		"hello info",
		"hello warn",
		"hello error",
	}
	zapLog, err := NewZapLogger(models.InfoLog, "./test.log")
	if err != nil {
		return
	}
	zapLog.Debug("info:", zap.String("s", s[0]))
	zapLog.Info("info:", zap.String("s", s[1]))
	zapLog.Warn("info:", zap.String("s", s[2]))
	zapLog.Error("info:", zap.String("s", s[3]))
}
