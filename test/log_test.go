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
	err = log.LoadLogFile(conf.Log)
	if err != nil {
		fmt.Println(err)
		return
	}
	//for i := 0; i < 20; i++ {
	//	go func() {
	//		for {
	log.Fatal("fatal", "", time.Now())
	log.Error("error", "", time.Now())
	log.Warn("warn", "", time.Now())
	log.Info("info", "", time.Now())
	log.Debug("debug", "", time.Now())
	log.Fatal("fatal", time.Now())
	log.Error("error", time.Now())
	log.Warn("warn", time.Now())
	log.Info("info", time.Now())
	log.Debug("debug", time.Now())
	log.Warn("warn", errors.New("xxxxx"), 123123123, time.Now())
	//		}
	//	}()
	//}
	//time.Sleep(time.Duration(20) * time.Second)
}
