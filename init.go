package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tabalt/pmon/process"
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

	writePid(config.PidFile)
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

// write pid
func writePid(file string) {
	err := process.WritePid(file, FileModeRW)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}
