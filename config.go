package main

import (
	"encoding/json"
	"io/ioutil"
)

type Process struct {
	Name     string `json:"name"`
	Enable   bool   `json:"enable"`
	User     string `json:"user"`
	Interval string `json:"interval"`
	PidFile  string `json:"pidfile"`

	Command string `json:"command"`
	StdOut  string `json:"stdout"`
	StdErr  string `json:"stderr"`
}

type Config struct {
	PidFile string `json:"pidfile"`
	LogFile string `json:"logfile"`

	ProcessList []*Process `json:"process"`
}

// init config data from file
func (c *Config) Init(file string) error {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, c)
	if err != nil {
		return err
	}

	return nil
}
