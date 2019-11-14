package log

import (
	"log"
	"os"
)

const (
	FatalLevel = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

const (
	FatalFlag = "[Fatal]"
	ErrorFlag = "[Error]"
	WarnFlag  = "[Warn]"
	InfoFlag  = "[Info]"
	DebugFlag = "[Debug]"
)

const (
	FlagTimeFmt  = ".20060102.15"
	ClearTimeFmt = ".20060102."
)

const (
	DefaultMaxSize    = 100 // max size(MB) for one log file
	DefaultExpireDays = 7   // reserved days for log files
	DefaultCheckDays  = 1
	DefaultLogLeval   = "debug"
)

var logger *log.Logger

var config *Config

var logLevel int

type Config struct {
	File       string
	MaxSize    int
	ExpireDays int
	CheckDays  int
	LogLevel   string
}

/**
 * load config file
 */
func LoadLogFile(conf Config) error {
	config = &Config{
		File:       conf.File,
		MaxSize:    conf.MaxSize,
		ExpireDays: conf.ExpireDays,
		CheckDays:  conf.CheckDays,
		LogLevel:   conf.LogLevel,
	}
	// init log level
	setLogLevel(config.LogLevel)
	// open log file
	file, err := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	// init logger
	logger = log.New(file, "", log.LstdFlags|log.Lmicroseconds)
	return nil
}

/**
 * write log
 */
func writeLog(flag interface{}, detail ...interface{}) {
	detail = append([]interface{}{flag}, detail...)
	logger.Println(detail...)
}

/**
 * write Fatal log
 */
func Fatal(detail ...interface{}) {
	if logLevel < FatalLevel {
		return
	}
	writeLog(FatalFlag, detail...)
}

/**
 * write Error log
 */
func Error(detail ...interface{}) {
	if logLevel < ErrorLevel {
		return
	}
	writeLog(ErrorFlag, detail...)
}

/**
 * write Warn log
 */
func Warn(detail ...interface{}) {
	if logLevel < WarnLevel {
		return
	}
	writeLog(WarnFlag, detail...)
}

/**
 * write Info log
 */
func Info(detail ...interface{}) {
	if logLevel < InfoLevel {
		return
	}
	writeLog(InfoFlag, detail...)
}

/**
 * write Debug log
 */
func Debug(detail ...interface{}) {
	if logLevel < DebugLevel {
		return
	}
	writeLog(DebugFlag, detail...)
}

/**
 * set log level
 */
func setLogLevel(level string) {
	switch level {
	case "fatal":
		logLevel = 0
	case "error":
		logLevel = 1
	case "warn":
		logLevel = 2
	case "info":
		logLevel = 3
	case "debug":
		logLevel = 4
	default:
		logLevel = 4
	}
}
