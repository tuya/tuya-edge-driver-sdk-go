/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

/*
Package logger provides a client for integration with the support-logging service. The client can also be configured
to write logs to a local file rather than sending them to a service.
*/
package logger

// Logging client for the Go implementation of edgexfoundry

import (
	"fmt"
	stdLog "log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/tuya/tuya-edge-driver-sdk-go/contracts/models"
	"github.com/tuya/tuya-edge-driver-sdk-go/internal/types"
)

// Colors
const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LoggingClient defines the interface for logging operations.
type LoggingClient interface {
	// SetLogLevel sets minimum severity log level. If a logging method is called with a lower level of severity than
	// what is set, it will result in no output.
	SetLogLevel(logLevel string) error
	// LogLevel returns the current log level setting
	LogLevel() string
	// Debug logs a message at the DEBUG severity level
	Debug(msg string, args ...interface{})
	// Error logs a message at the ERROR severity level
	Error(msg string, args ...interface{})
	// Info logs a message at the INFO severity level
	Info(msg string, args ...interface{})
	// Trace logs a message at the TRACE severity level
	Trace(msg string, args ...interface{})
	// Warn logs a message at the WARN severity level
	Warn(msg string, args ...interface{})
	// Debugf logs a formatted message at the DEBUG severity level
	Debugf(msg string, args ...interface{})
	// Errorf logs a formatted message at the ERROR severity level
	Errorf(msg string, args ...interface{})
	// Infof logs a formatted message at the INFO severity level
	Infof(msg string, args ...interface{})
	// Tracef logs a formatted message at the TRACE severity level
	Tracef(msg string, args ...interface{})
	// Warnf logs a formatted message at the WARN severity level
	Warnf(msg string, args ...interface{})
}

type edgeXLogger struct {
	owningServiceName string
	logLevel          *string
	//rootLogger        log.Logger
	//levelLoggers      map[string]log.Logger
	rootLogger   *zap.Logger
	levelLoggers map[string]*zap.Logger
}

const (
	LogPathEnvName = "LOG_PATH"
	LogLevel       = "LOG_LEVEL"
)

var (
	levels = map[string]zapcore.Level{
		models.DebugLog: zap.DebugLevel,
		models.InfoLog:  zap.InfoLevel,
		models.WarnLog:  zap.WarnLevel,
		models.ErrorLog: zap.ErrorLevel,
	}
)

type LogMessage struct {
	Time        string        `json:"time"`
	ServiceName string        `json:"service_name"`
	Caller      string        `json:"caller"`
	Message     []interface{} `json:"message"`
}

func NewZapLogger(logLevelStr, logPath string) (zapLog *zap.Logger, err error) {
	stdLog.SetFlags(stdLog.LstdFlags | stdLog.Llongfile)
	//encoderConfig := zap.NewProductionEncoderConfig()
	// 选择自定义日志样式
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		//StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var logLevel zapcore.Level
	var exists bool
	logLevel, exists = levels[logLevelStr]
	if !exists {
		logLevel = zap.DebugLevel
	}

	if logPath != "" {
		logDir := path.Dir(logPath)
		if _, err = os.Stat(logDir); os.IsNotExist(err) {
			stdLog.Fatal("ERROR 日志目录 ", logDir, " 不存在")
			return
		}
		// 打印到文件，自动分裂
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath,
			MaxSize:    64, // megabytes
			MaxBackups: 10,
			MaxAge:     28, // days
		})
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			w,
			zap.NewAtomicLevelAt(logLevel),
		)
		zapLog = zap.New(core, zap.AddCaller())
	} else {
		// 打印到控制台
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(logLevel)
		cfg.Encoding = "console"
		cfg.EncoderConfig = encoderConfig
		//cfg.DisableStacktrace = true
		zapLog, err = cfg.Build()
		if err != nil {
			stdLog.Fatal("ERROR ", err)
			return
		}
	}
	return
}

// NewClient creates an instance of LoggingClient
func NewClient(owningServiceName string, logLevel string) LoggingClient {
	// 从Env环境变量获取日志等级
	logLevel = os.Getenv(LogLevel)

	if !isValidLogLevel(logLevel) {
		logLevel = models.DebugLog
	}

	// Set up logging client
	lc := edgeXLogger{
		owningServiceName: owningServiceName,
		logLevel:          &logLevel,
	}

	//lc.rootLogger = log.NewLogfmtLogger(os.Stdout)
	//lc.rootLogger = log.WithPrefix(
	//	lc.rootLogger,
	//	"ts",
	//	log.DefaultTimestampUTC,
	//	"app",
	//	owningServiceName,
	//	"source",
	//	log.Caller(5))
	//
	//// Set up the loggers
	//lc.levelLoggers = map[string]log.Logger{}

	var err error
	lc.rootLogger, err = NewZapLogger(logLevel, os.Getenv(LogPathEnvName))
	if err != nil {
		return nil
	}
	lc.levelLoggers = make(map[string]*zap.Logger)
	for _, level := range logLevels() {
		lc.levelLoggers[level] = lc.rootLogger
		//lc.levelLoggers[level] = log.WithPrefix(lc.rootLogger, "level", level)
	}

	return lc
}

// LogLevels returns an array of the possible log levels in order from most to least verbose.
func logLevels() []string {
	return []string{
		models.TraceLog,
		models.DebugLog,
		models.InfoLog,
		models.WarnLog,
		models.ErrorLog,
	}
}

func isValidLogLevel(l string) bool {
	for _, name := range logLevels() {
		if name == l {
			return true
		}
	}
	return false
}

func (lc edgeXLogger) log(logLevel string, formatted bool, msg string, args ...interface{}) {
	// Check minimum log level
	for _, name := range logLevels() {
		if name == *lc.logLevel {
			break
		}
		if name == logLevel {
			return
		}
	}

	if args == nil {
		args = []interface{}{msg}
	} else if formatted {
		args = []interface{}{fmt.Sprintf(msg, args...)}
	} else {
		//if len(args)%2 == 1 {
		//	// add an empty string to keep k/v pairs correct
		//	args = append(args, "")
		//}
		if len(msg) > 0 {
			args = append(args, msg)
		}
	}
	argData := make([]string, 0)
	for _, arg := range args {
		argData = append(argData, fmt.Sprintf("%+v", arg))
	}
	msg = strings.Join(argData, ",")

	//err := lc.levelLoggers[logLevel].Log(args...)
	//if err != nil {
	//	stdLog.Fatal(err.Error())
	//	return
	//}

	//_, file, line, _ := runtime.Caller(2)
	//idx := strings.LastIndexByte(file, '/')
	//caller := file[idx+1:] + ":" + strconv.Itoa(line)

	// 日志样式设置
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	paths := strings.Split(funcName, "/")
	pack := strings.Split(funcName[strings.LastIndexByte(funcName, '/')+1:], ".")[0]
	paths = append(paths[:len(paths)-1], pack, strings.Split(file[strings.LastIndexByte(file, '/')+1:], ".")[0])
	funcPath := strings.Join(paths, "/")
	caller := funcPath + ".go:" + strconv.Itoa(line)
	now := time.Now().Format("2006/01/02 15:04:05")
	debugStr := Yellow + "[debug] " + Reset
	debugCaller := fmt.Sprintf(" %v%v%v ", Yellow, caller, Reset)
	infoStr := Green + "[info] " + Reset
	infoCaller := fmt.Sprintf(" %v%v%v ", Green, caller, Reset)
	warnStr := Magenta + "[warn] " + Reset
	warnCaller := fmt.Sprintf(" %v%v%v ", Magenta, caller, Reset)
	errStr := Red + "[error] " + Reset
	errCaller := fmt.Sprintf(" %v%v%v ", Red, caller, Reset)
	var message = fmt.Sprintf("%v %v[%v]%v", now, Blue, lc.owningServiceName, Reset)

	// 日志输出
	switch logLevel {
	case models.DebugLog:
		lc.levelLoggers[logLevel].Debug(debugStr + message + debugCaller + msg)
	case models.InfoLog:
		lc.levelLoggers[logLevel].Info(infoStr + message + infoCaller + msg)
	case models.WarnLog:
		lc.levelLoggers[logLevel].Warn(warnStr + message + warnCaller + msg)
	case models.ErrorLog:
		lc.levelLoggers[logLevel].Error(errStr + message + errCaller + msg)
	}
}

func (lc edgeXLogger) SetLogLevel(logLevel string) error {
	if isValidLogLevel(logLevel) == true {
		*lc.logLevel = logLevel

		return nil
	}

	return types.ErrNotFound{}
}

func (lc edgeXLogger) LogLevel() string {
	if lc.logLevel == nil {
		return ""
	}
	return *lc.logLevel
}

func (lc edgeXLogger) Info(msg string, args ...interface{}) {
	lc.log(models.InfoLog, false, msg, args...)
}

func (lc edgeXLogger) Trace(msg string, args ...interface{}) {
	lc.log(models.TraceLog, false, msg, args...)
}

func (lc edgeXLogger) Debug(msg string, args ...interface{}) {
	lc.log(models.DebugLog, false, msg, args...)
}

func (lc edgeXLogger) Warn(msg string, args ...interface{}) {
	lc.log(models.WarnLog, false, msg, args...)
}

func (lc edgeXLogger) Error(msg string, args ...interface{}) {
	lc.log(models.ErrorLog, false, msg, args...)
}

func (lc edgeXLogger) Infof(msg string, args ...interface{}) {
	lc.log(models.InfoLog, true, msg, args...)
}

func (lc edgeXLogger) Tracef(msg string, args ...interface{}) {
	lc.log(models.TraceLog, true, msg, args...)
}

func (lc edgeXLogger) Debugf(msg string, args ...interface{}) {
	lc.log(models.DebugLog, true, msg, args...)
}

func (lc edgeXLogger) Warnf(msg string, args ...interface{}) {
	lc.log(models.WarnLog, true, msg, args...)
}

func (lc edgeXLogger) Errorf(msg string, args ...interface{}) {
	lc.log(models.ErrorLog, true, msg, args...)
}

// Build the log entry object
func (lc edgeXLogger) buildLogEntry(logLevel string, msg string, args ...interface{}) models.LogEntry {
	res := models.LogEntry{}
	res.Level = logLevel
	res.Message = msg
	res.Args = args
	res.OriginService = lc.owningServiceName

	return res
}
