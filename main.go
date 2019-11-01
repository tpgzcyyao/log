package main

import (
	"fmt"
	"github.com/tpgzcyyao/config"
	"os"
)

func main() {
	fmt.Println("log")
	res, err := (new(config.Config)).LoadFile("test.conf")
	if err != nil {
		fmt.Println("LoadFile error.")
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("%v", res))
}
