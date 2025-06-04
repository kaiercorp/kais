package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *logrus.Logger

type myFormatter struct {
	logrus.TextFormatter
}

func (f *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s][%s] %s\n",
		strings.ToUpper(entry.Level.String()),
		entry.Time.Format(f.TimestampFormat),
		entry.Message)), nil
}

func initCommon(loglevel logrus.Level, log_file string) {
	lum := &lumberjack.Logger{
		Filename:   log_file,
		MaxSize:    100,
		MaxBackups: 1000,
		MaxAge:     90,
		LocalTime:  true,
		Compress:   false,
	}

	logger = &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout, lum), //os.Stdout,//
		Level: loglevel,
	}

	logger.SetFormatter(&myFormatter{logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	},
	})
}

func InitDebug(log_file string) {
	initCommon(logrus.DebugLevel, log_file)
}

func InitInfo(log_file string) {
	initCommon(logrus.InfoLevel, log_file)
}

func InitLogger(log_level string, log_file string) {
	if log_level == "DEBUG" {
		InitDebug(log_file)
	} else {
		InitInfo(log_file)
	}
}

func getCaller(depth int) string {
	_, file, no, ok := runtime.Caller(depth)
	if ok {
		_, sfile := filepath.Split(file)
		return fmt.Sprintf("%s:%d", sfile, no)
	}
	return ""
}

func setMessageWithCaller(message string) string {
	caller := getCaller(3)
	if caller == "" {
		return message
	}
	return fmt.Sprintf("(%s) %+v", caller, message)
}

func Debug(args ...interface{}) {
	msg := make_message(args...)
	logger.Debug(setMessageWithCaller(msg))
}

func Info(args ...interface{}) {
	msg := make_message(args...)
	logger.Info(setMessageWithCaller(msg))
}

func Warn(args ...interface{}) {
	msg := make_message(args...)
	logger.Warn(setMessageWithCaller(msg))
}

func Error(args ...interface{}) {
	msg := make_message(args...)
	logger.Error(setMessageWithCaller(msg))
}

func Panic(args ...interface{}) {
	msg := make_message(args...)
	logger.Panic(setMessageWithCaller(msg))
}

func Fatal(args ...interface{}) {
	msg := make_message(args...)
	logger.Fatal(setMessageWithCaller(msg))
}

func make_message(args ...interface{}) string {
	format := fmt.Sprintf("%v", args[0])
	if strings.Contains(format, "%s") || strings.Contains(format, "%d") {
		argv := args[1:]
		return fmt.Sprintf(format, argv...)
	}

	return fmt.Sprint(args...)
}

func ReportError(code string, message string, e error) {
	msg := ""
	if e != nil {
		msg = fmt.Sprintf("[%s] %s: %s", code, message, e.Error())
	} else {
		msg = fmt.Sprintf("[%s] %s", code, message)
	}
	caller := getCaller(3)
	logger.Error(fmt.Sprintf("(%s) %s", caller, msg))
}
