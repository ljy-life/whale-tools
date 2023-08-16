package config

import (
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

const (
	LOG_NAME   = "./debug.log"
	LOG_SIZE   = 50
	LOG_BACKUP = 10
	LOG_DATE   = 7
)

var Log = InitLogger()

func InitLogger() *logrus.Logger {
	// 对日志文件进行大小轮转等设置
	logconf := &lumberjack.Logger{
		Filename:   LOG_NAME,
		MaxSize:    LOG_SIZE,
		MaxAge:     LOG_DATE,
		MaxBackups: LOG_BACKUP,
		Compress:   true,
	}
	logger := logrus.New()

	// 设置 textformatter 的文本格式
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置 jsonformatter 的文本格式
	//logrus.SetFormatter(&logrus.JSONFormatter{
	//	TimestampFormat:   "2006-01-02 15:04:05",
	//	PrettyPrint:       true,
	//})

	// 将日志同时输出至桌面和日志文件中
	fileAndStdoutWriter := io.MultiWriter(os.Stdout, logconf)
	logger.SetOutput(fileAndStdoutWriter)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}
