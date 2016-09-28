package main

import (
	"fmt"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/tabalt/pmon/process"
)

const (
	DefaultMonitorInterval = 10 * time.Second
	DefaultStartWait       = 10 * time.Second
)

func main() {
	logger.Println("pmon started")

	complete := make(chan int)

	processCount := 0
	for _, ps := range config.ProcessList {
		if !ps.Enable || ps.PidFile == "" {
			continue
		}

		processCount++
		go monitorProcess(ps, complete)
	}

	for i := 0; i < processCount; i++ {
		<-complete
	}

	logger.Println("pmon shutting down")
}

// monitor a process
func monitorProcess(ps *Process, complete chan int) {
	for {
		// monitor process
		logger.Println("monitor process " + ps.Name + " by pid file " + ps.PidFile)

		pid, err := process.ReadPid(ps.PidFile)
		if err != nil {
			logger.Println("failed to get pid, error: " + err.Error())
			startProcess(ps)
			continue
		}

		running, err := process.IsRunning(pid)
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

		// sleep a while
		d, intervalError := time.ParseDuration(ps.Interval)
		if intervalError != nil {
			d = DefaultMonitorInterval
		}
		time.Sleep(d)

	}
	complete <- 1
}

// try to start process
func startProcess(ps *Process) {
	shell := fmt.Sprintf("nohup %s 1%s 2%s &", ps.Command, ps.StdOut, ps.StdErr)

	logger.Println(fmt.Sprintf("try to start %s, exec bash command: %s", ps.Name, shell))

	_, err := runShell(shell, ps.User)
	if err != nil {
		logger.Println(fmt.Sprintf("start %s failed, error: %v", ps.Name, err))
		return
	}

	// wait process start
	d, err := time.ParseDuration(ps.StartWait)
	if err != nil {
		d = DefaultStartWait
	}
	time.Sleep(d)

}

// run shell
func runShell(shell string, userName string) ([]byte, error) {
	uc, _ := user.Current()
	if userName != "" && userName != uc.Username {
		if uc.Username == "root" {
			shell = fmt.Sprintf(
				"/sbin/runuser %s -c \"%s\"",
				userName,
				strings.Replace(shell, "\"", "\\\"", -1),
			)
		} else {
			shell = fmt.Sprintf(
				"sudo -u %s %s",
				userName,
				shell,
			)
		}
	}

	cmd := exec.Command("/bin/bash", "-c", shell)
	return cmd.CombinedOutput()
}
