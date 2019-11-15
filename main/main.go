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
	err := (new(config.Config)).LoadConfig("/export/liyang/gopath/src/log/test/test.conf", conf)
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
