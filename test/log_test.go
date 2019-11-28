package test

import (
    "errors"
    "fmt"
    "os"
    "testing"
    "time"

    "github.com/tpgzcyyao/config"
    "github.com/tpgzcyyao/log"
)

type Config struct {
    Log log.Config
}

func TestLog(t *testing.T) {
    conf := new(Config)
    dir, _ := os.Getwd()
    err := (new(config.Config)).LoadConfig(dir+"/test.conf", conf)
    // fileName should be absolute path in config file
    conf.Log.FileName = dir + "/../log/test.log"
    fmt.Println(conf)
    if err != nil {
        fmt.Println(err)
        return
    }
    err = log.LoadLogConfig(conf.Log)
    if err != nil {
        fmt.Println(err)
        return
    }
    //for i := 0; i < 20; i++ {
    //    go func() {
    //        for {
    log.Fatal("This is a fatal.", "1", time.Now())
    log.Error("This is an error.", "xxxx", time.Now())
    log.Warn("This is a warn.", 123, time.Now())
    log.Info("This is an info.", "", time.Now())
    log.Debug("This is a debug.", "xxx", time.Now())
    log.Fatal("fatal", time.Now())
    log.Error("error", time.Now())
    log.Warn("warn", time.Now())
    log.Info("info", time.Now())
    log.Debug("debug", time.Now())
    log.Warn("warn", errors.New("xxxxx"), 123123123, time.Now())
    log.Fatal("This is a fatal.", 123, "xxx")
    log.Error("This is an error.")
    log.Warn("This is a warn.")
    log.Info("This is a info.")
    log.Debug("This is a debug.")
    //        }
    //    }()
    //}
    //time.Sleep(time.Duration(300) * time.Second)
}
