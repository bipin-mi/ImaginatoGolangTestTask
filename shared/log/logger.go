package log

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"bytes"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Data is
type Data struct {
	IPAddress string ``
	Session   string ``
	ActorID   string ``
	ActorType string ``
}

// ILogger is
type ILogger interface {
	Debug(data interface{}, description string, args ...interface{})
	Info(data interface{}, description string, args ...interface{})
	Warn(data interface{}, description string, args ...interface{})
	Error(data interface{}, description string, args ...interface{})
	Fatal(data interface{}, description string, args ...interface{})
	Panic(data interface{}, description string, args ...interface{})
}

type logger struct {
	appName    string
	appVersion string

	filePath string
	level    logrus.Level
	maxAge   time.Duration

	theLogger *logrus.Logger
}

var defaultLogger logger
var defaultLoggerOnce sync.Once

//createLogger is
func createLogger() {
	formatter := nested.Formatter{
		NoColors:        true,
		HideKeys:        true,
		TimestampFormat: "0102 150405.000",
		FieldsOrder:     []string{"func"},
	}

	defaultLogger.theLogger = logrus.New()
	defaultLogger.theLogger.SetLevel(defaultLogger.level)
	defaultLogger.theLogger.SetFormatter(&formatter)

	filename := defaultLogger.appName + ".%Y%m%d.log"
	if len(defaultLogger.appVersion) > 0 {
		filename = defaultLogger.appName + "-" + defaultLogger.appVersion + ".%Y%m%d.log"
	}

	writer, _ := rotateLogs.New(
		filepath.Join(defaultLogger.filePath, filename),
		rotateLogs.WithMaxAge(defaultLogger.maxAge),
		rotateLogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	defaultLogger.theLogger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.WarnLevel:  writer,
			logrus.ErrorLevel: writer,
			logrus.DebugLevel: writer,
		},
		defaultLogger.theLogger.Formatter,
	))
}

//Init initialize the logger
func Init(appName, filePath string, logLevel logrus.Level, maxAge time.Duration) {
	defaultLogger.appName = appName
	defaultLogger.filePath = filePath
	defaultLogger.level = logLevel
	defaultLogger.maxAge = maxAge
}

// GetLog is to retrieve the logs
func GetLog() ILogger {
	defaultLoggerOnce.Do(createLogger)
	return &defaultLogger
}

//getLogEntry is
func (selfLogger *logger) getLogEntry(extraInfo interface{}) *logrus.Entry {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()

	var buffer bytes.Buffer

	buffer.WriteString("fn:")

	x := strings.LastIndex(funcName, "/")
	buffer.WriteString(funcName[x+1:])

	if extraInfo == nil {
		return selfLogger.theLogger.WithField("info", buffer.String())
	}

	data, ok := extraInfo.(Data)
	if !ok {
		return selfLogger.theLogger.WithField("info", buffer.String())
	}

	if data.IPAddress != "" {
		buffer.WriteString("|ip:")
		buffer.WriteString(data.IPAddress)
	}

	if data.Session != "" {
		buffer.WriteString("|ss:")
		buffer.WriteString(data.Session)
	}

	if data.ActorID != "" {
		buffer.WriteString("|id:")
		buffer.WriteString(data.ActorID)
	}

	if data.ActorType != "" {
		buffer.WriteString("|tp:")
		buffer.WriteString(data.ActorType)
	}

	return selfLogger.theLogger.WithField("info", buffer.String())
}

// Debug is
func (selfLogger *logger) Debug(data interface{}, description string, args ...interface{}) {
	selfLogger.getLogEntry(data).Debugf(description+"\n", args...)
}

// Info is
func (selfLogger *logger) Info(data interface{}, description string, args ...interface{}) {
	selfLogger.getLogEntry(data).Infof(description+"\n", args...)
}

// Warn is
func (selfLogger *logger) Warn(data interface{}, description string, args ...interface{}) {
	selfLogger.getLogEntry(data).Warnf(description+"\n", args...)
}

// Error is
func (selfLogger *logger) Error(data interface{}, description string, args ...interface{}) {
	selfLogger.getLogEntry(data).Errorf(description+"\n", args...)
}

// Fatal is
func (selfLogger *logger) Fatal(data interface{}, description string, args ...interface{}) {
	selfLogger.getLogEntry(data).Fatalf(description+"\n", args...)
}

// Panic is
func (selfLogger *logger) Panic(data interface{}, description string, args ...interface{}) {
	selfLogger.getLogEntry(data).Panicf(description+"\n", args...)
}
