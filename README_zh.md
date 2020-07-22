# log
[English Document](./README.md)
## 一、介绍
这是一个用golang语言实现的记录程序运行日志的包。
## 二、特性
- 支持Fatal、Error、Warn、Info、Debug五种级别的日志。
- 根据文件大小和日期进行分割。
- 根据日志保留时间配置进行过期日志清除。
- 日志写入方法支持传入多个不确定类型的参数。
## 三、使用方法
- 下载包
```
go get github.com/tpgzcyyao/log
```
- 导入包
```
import "github.com/typzcyyao/log"
```
- 代码示例
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
- 输出日志示例
```
2019/11/15 15:57:17.351581 [Fatal] This is a fatal. 123 xxx
2019/11/15 15:57:17.351588 [Error] This is an error.
2019/11/15 15:57:17.351595 [Warn] This is a warn.
2019/11/15 15:57:17.351601 [Info] This is an info.
```
## 四、配置说明
- log.Config.FileName
日志的文件名，完整的绝对路径，必须配置。
- log.Config.MaxSize
每个日志文件的最大大小，单位为MB，默认为100。
- log.Config.ExpireDays
日志的保留天数，默认全保留。
- log.Config.LogLevel
打印日志的最低等级，值为fatal、error、warn、info、debug，默认为debug。
- log.Config.TotalSize
所有日志文件的总大小，单位为MB。
## 五、关联配置文件
- 代码
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
- 配置文件 /export/config/test.conf
```
[log]
file_name = /export/log/test.log
max_size = 200
expire_days = 7
log_level = info
total_size = 1024
```
