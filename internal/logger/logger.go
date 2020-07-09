package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

type Logger struct {
	Log *logrus.Logger
}

func ConfigureLogger() *Logger {
	logger := new(Logger)
	logger.Log = logrus.New()
	logger.Log.Formatter = &logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano}
	logger.Log.Level = logrus.InfoLevel
	logger.Log.Out = os.Stdout
	return logger
}

func (l Logger) LogWithFields(req *http.Request, level string, fields map[string]interface{}) {
	msg := ""
	_, ok := fields["msg"]
	if ok {
		msg = fmt.Sprintf("%v", fields["msg"])
		delete(fields, "msg")
	}

	if req != nil {
		fields["path"] = req.URL.Path
		fields["header"] = req.Header
		fields["reqMethod"] = req.Method
	}

	logEntry := l.Log.WithFields(fields)

	switch strings.ToLower(level) {
	case "panic":
		logEntry.Panic(msg)
		return
	case "fatal":
		logEntry.Fatal(msg)
		return
	case "error":
		logEntry.Error(msg)
		return
	case "warn", "warning":
		logEntry.Warn(msg)
		return
	case "info":
		logEntry.Info(msg)
		return
	case "debug":
		logEntry.Debug(msg)
		return
	}
	logEntry.Info(msg) // Default level
}
