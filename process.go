package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

const (
	DefaultMonitorInterval = 10 * time.Second
	DefaultStartWait       = 10 * time.Second

	ProcessStatIndex = 2

	ProcessStatRunning         = "R"
	ProcessStatSleeping        = "S"
	ProcessStatStoped          = "T"
	ProcessStatZombie          = "Z"
	ProcessStatUninterruptible = "D"
)

// monitor a process
func monitorProcess(ps *Process, complete chan int) {
	for {

		// sleep a while
		d, intervalError := time.ParseDuration(ps.Interval)
		if intervalError != nil {
			d = DefaultMonitorInterval
		}
		time.Sleep(d)

		// monitor process
		logger.Println("monitor process " + ps.Name + " by pid file " + ps.PidFile)

		pid, err := getPidFromFile(ps.PidFile)
		if err != nil {
			logger.Println("failed to get pid, error: " + err.Error())
			startProcess(ps)
			continue
		}

		running, err := isProcessRunning(pid)
		if err != nil {
			logger.Println("failed to get process stat, error: " + err.Error())
			startProcess(ps)
			continue
		}

		if !running {
			logger.Println("process with pid " + pid + " not running")
			startProcess(ps)
			continue
		}

		logger.Println("process with pid " + pid + " is running")

	}
	complete <- 1
}

// get pid from pid file
func getPidFromFile(file string) (string, error) {
	pidBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	pid := strings.TrimSpace(string(pidBytes))
	if pid == "" {
		return "", fmt.Errorf("pid file is empty")
	}
	return pid, nil
}

// check process running or not by pid
func isProcessRunning(pid string) (bool, error) {
	stat, err := getProcessStatByPid(pid)
	if err != nil {
		return false, err
	}

	if stat == ProcessStatStoped || stat == ProcessStatZombie {
		return false, nil
	}
	return true, nil
}

// read /proc/$pid/stat to get process stat
func getProcessStatByPid(pid string) (string, error) {
	file := "/proc/" + pid + "/stat"
	statBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	stat := strings.TrimSpace(string(statBytes))
	statList := strings.Split(stat, " ")

	if len(statList) < (ProcessStatIndex + 1) {
		return "", fmt.Errorf("stat file is empty")
	}

	return strings.TrimSpace(statList[ProcessStatIndex]), nil
}

// try to start process
func startProcess(ps *Process) (string, error) {

	shell := fmt.Sprintf("%s 1%s 2%s &", ps.Command, ps.StdOut, ps.StdErr)

	logger.Println(fmt.Sprintf(
		"try to start %s, exec bash command: %s",
		ps.Name,
		shell,
	))

	cmd := exec.Command("/bin/bash", "-c", shell)
	cmd.Start()

	// wait process start
	d, err := time.ParseDuration(ps.StartWait)
	if err != nil {
		d = DefaultStartWait
	}
	time.Sleep(d)

	return "", nil
}
