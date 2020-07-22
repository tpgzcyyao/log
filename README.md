# log
[中文文档](./README_zh.md)
## I. Introduction
This is a golang package using for writing log.
## II. Features
- Support Fatal, Error, Warn, Info, Debug five log levels.
- Split log file by size and date.
- Clear expired log files.
- Log write method supports passing in multiple parameters of uncertain types.
## III. Description
- Download the package.
```
go get github.com/tpgzcyyao/log
```
- Import the package.
```
import "github.com/typzcyyao/log"
```
- Example
```
logConfig := log.Config{
	FileName: "/export/log/test.log",
	MaxSize: 200,
	ExpireDays: 7,
	LogLevel: "info",
	TotalSize: 1024,
}
err := LoadLogConfig(config)
if err != nil {
	panic(err)
}
log.Fatal("This is a fatal.", 123, "xxx")
log.Error("This is an error.")
log.Warn("This is a warn.")
log.Info("This is an info.")
log.Debug("This is a debug.")
```
- Output in the log file
```
2019/11/15 15:57:17.351581 [Fatal] This is a fatal. 123 xxx
2019/11/15 15:57:17.351588 [Error] This is an error.
2019/11/15 15:57:17.351595 [Warn] This is a warn.
2019/11/15 15:57:17.351601 [Info] This is an info.
```
## IV. Config Instructions
- log.Config.FileName
FileName represents the path for log file. It must be completed absolute path.
- log.Config.MaxSize
MaxSize represents the max size for each log file and the unit is MB. Default max size is 100.
- log.Config.ExpireDays
ExpireDays represents the number of days to keep log. Default config is to keep all log files.
- log.Config.LogLevel
LogLevel represents the lowest level for print log. The value may be fatal, error, warn, info, debug. Default log level is debug.
- log.Config.TotalSize
TotalSize represents the total size for all log files and the unit is MB.
## V. Using Config File
- Code
```
package main

import (
	"fmt"

	"github.com/tpgzcyyao/config"
	"github.com/tpgzcyyao/log"
)

type Config struct {
	Log log.Config
}

func main() {
	conf := new(Config)
	err := (new(config.Config)).LoadConfig("/export/config/test.conf", conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = log.LoadLogConfig(conf.Log)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Fatal("This is a fatal.", 123, "xxx")
	log.Error("This is an error.")
	log.Warn("This is a warn.")
	log.Info("This is a info.")
	log.Debug("This is a debug.")
}
```
- config file: /export/config/test.conf
```
[log]
file_name = /export/log/test.log
max_size = 200
expire_days = 7
log_level = info
total_size = 1024
```
