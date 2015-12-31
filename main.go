package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const (
	FileModeRW os.FileMode = 0666
)

var (
	configFile string

	config *Config
	logger *log.Logger
)

func init() {
	initFlag()
	initConfig()
	initLogger(config.LogFile)

	savePid(config.PidFile)
}

func main() {
	logger.Println("pmon started")

	complete := make(chan int)
	for _, ps := range config.ProcessList {
		if !ps.Enable || ps.PidFile == "" {
			continue
		}

		go monitorProcess(ps, complete)
		<-complete
	}

	logger.Println("pmon shutting down")
}

// init flag
func initFlag() {
	flag.StringVar(&configFile, "c", "./pmon.json", "config file for pmon")
	flag.Parse()
}

// init config
func initConfig() {
	config = &Config{}
	err := config.Init(configFile)
	if err != nil {
		fmt.Printf("init config failed. error: %s\n", err.Error())
		os.Exit(1)
	}
}

// init logger
func initLogger(file string) {
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, FileModeRW)
	if err != nil {
		fmt.Printf("init log failed. error: %s\n", err.Error())
		os.Exit(1)
	}
	logger = log.New(logFile, "", log.Ldate|log.Ltime)
}

// save pid
func savePid(file string) {
	err := ioutil.WriteFile(file, []byte(strconv.Itoa(os.Getpid())), FileModeRW)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}
