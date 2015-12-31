package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func main() {
	savePid("./testdaemon/tmp/testdaemon.pid")
	for {
		fmt.Println("hi~")

		d, intervalError := time.ParseDuration("5s")
		if intervalError != nil {
			os.Exit(1)
		}
		time.Sleep(d)
	}
}

// save pid
func savePid(file string) {
	err := ioutil.WriteFile(file, []byte(strconv.Itoa(os.Getpid())), 0666)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
