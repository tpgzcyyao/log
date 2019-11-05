package log

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	FatalFlag    = "Fatal"
	ErrorFlag    = "Error"
	WarnFlag     = "Warn"
	InfoFlag     = "Info"
	DebugFlag    = "Debug"
	FlagTimeFmt  = ".20060102.15"
	ClearTimeFmt = ".20060102."
)

const (
	DefaultLenChannel = 1024
	DefaultExpireDays = 7
	DefaultCheckDays  = 7 // 删除过期时间点前7天的日志，再之前的不处理
)

type Config struct {
	File       string
	ChanLen    int
	ExpireDays int
}

type Content struct {
	Flag   string
	Msg    string
	Detail string
}

var logger *log.Logger

var chLog chan *Content

var config *Config

// 文件分割时间标记
var flagTime string

// 文件清理时间标记
var clearTime string

// 文件分割锁
var splitLock bool

/**
 * 加载日志配置文件
 */
func LoadLogFile(conf Config) error {
	config = &Config{
		File:       conf.File,
		ChanLen:    conf.ChanLen,
		ExpireDays: conf.ExpireDays,
	}
	// 打开日志文件
	file, err := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("打开日志文件失败！")
	}
	// 设置Logger
	logger = log.New(file, "", log.LstdFlags|log.Lmicroseconds)
	// 初始化日志通道
	chanLenTmp := config.ChanLen
	if chanLenTmp <= 0 {
		chanLenTmp = DefaultLenChannel
	}
	chLog = make(chan *Content, chanLenTmp)
	// 初始化锁状态
	splitLock = false
	// 初始化文件时间标记
	flagTime = time.Now().Format(FlagTimeFmt)
	clearTime = ""
	// 写入文件
	go writeToFile()
	// 日志split检查
	go splitCheck()
	return nil
}

/**
 * send logs to channel
 */
func printToChan(flag string, msg string, detail string) {
	chLog <- &Content{Flag: flag, Msg: msg, Detail: detail}
}

/**
 * write Fatal log
 */
func Fatal(msg string, detail string) {
	printToChan(FatalFlag, msg, detail)
}

/**
 * write Error log
 */
func Error(msg string, detail string) {
	printToChan(ErrorFlag, msg, detail)
}

/**
 * write Warn log
 */
func Warn(msg string, detail string) {
	printToChan(WarnFlag, msg, detail)
}

/**
 * write Info log
 */
func Info(msg string, detail string) {
	printToChan(InfoFlag, msg, detail)
}

/**
 * write Debug log
 */
func Debug(msg string, detail string) {
	printToChan(DebugFlag, msg, detail)
}

/**
 * 日志写入文件
 *
 * 从管道中读取日志信息，写入到文件
 */
func writeToFile() {
	for {
		if splitLock {
			// 移动日志文件
			fileBak := config.File + flagTime
			_ = os.Rename(config.File, fileBak)
			// 重新设置日志输出文件
			file, _ := os.OpenFile(config.File, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			logger.SetOutput(file)
			// 释放分割锁
			splitLock = false
			// 日志过期清理
			go clearExpiredLog()
			continue
		}
		content := <-chLog
		logger.Println(content.Flag, content.Msg, content.Detail)
	}
}

/**
 * 判断日志分割的时机
 *
 * 整点切分
 */
func splitCheck() {
	for {
		// 有锁不操作
		if splitLock {
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}
		// 判断当前时间
		nowTime := time.Now().Format(FlagTimeFmt)
		if nowTime == flagTime {
			time.Sleep(time.Duration(1) * time.Second)
			continue
		}
		// 更新分割时间标记
		flagTime = nowTime
		// 锁定分割锁
		splitLock = true
	}
}

/**
 * 日志过期清理
 */
func clearExpiredLog() {
	// 判断日志清理标记
	timeNow := time.Now().Format(ClearTimeFmt)
	if timeNow == clearTime {
		return
	}
	// 判断保留日志天数
	expireDays := config.ExpireDays
	if expireDays <= 0 {
		expireDays = DefaultExpireDays
	}
	// 遍历需要清理的文件
	for i := expireDays + 1; i < expireDays+DefaultCheckDays+1; i++ {
		day := time.Now().Add(-time.Hour * time.Duration(24*i)).Format(ClearTimeFmt)
		for j := 0; j < 24; j++ {
			time := strconv.Itoa(j)
			if len(time) < 2 {
				time = "0" + time
			}
			fileBak := config.File + day + time
			os.Remove(fileBak)
		}
	}
	// 更新清理时间标记
	clearTime = timeNow
}
