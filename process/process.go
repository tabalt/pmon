package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	ProcStatFieldIndex = 2

	StatRunning         = "R"
	StatSleeping        = "S"
	StatStoped          = "T"
	StatZombie          = "Z"
	StatUninterruptible = "D"
)

// read pid from a file
func ReadPid(file string) (string, error) {
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

// write pid to a file
func WritePid(file string, mode os.FileMode) error {
	return ioutil.WriteFile(file, []byte(strconv.Itoa(os.Getpid())), mode)
}

// check a process running or not by pid
func IsRunning(pid string) (bool, error) {
	stat, err := GetStatByPid(pid)
	if err != nil {
		return false, err
	}

	if stat == StatStoped || stat == StatZombie {
		return false, nil
	}
	return true, nil
}

// read /proc/$pid/stat to get process stat
func GetStatByPid(pid string) (string, error) {
	file := "/proc/" + pid + "/stat"
	statBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	stat := strings.TrimSpace(string(statBytes))
	statList := strings.Split(stat, " ")

	if len(statList) < (ProcStatFieldIndex + 1) {
		return "", fmt.Errorf("stat file is empty")
	}

	return strings.TrimSpace(statList[ProcStatFieldIndex]), nil
}
