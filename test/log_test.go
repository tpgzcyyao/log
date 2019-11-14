package test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/tpgzcyyao/log"
)

func TestLog(t *testing.T) {
	config := log.Config{
		File:     "/export/liyang/gopath/src/log/log/test.log",
		LogLevel: "debug",
	}
	err := log.LoadLogFile(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 20; i++ {
		go func() {
			for {
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
			}
		}()
	}
	time.Sleep(time.Duration(20) * time.Second)
}
