package util

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/v3/process"
)

func GetExecDirectory() string {
	file, err := os.Getwd()
	if err == nil {
		return file + "/"
	}
	return ""
}

func CheckProcessExist(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = p.Signal(syscall.Signal(0))
	return err == nil
}

func CheckProcessExistOptionalName(pid int, name string) (bool, error) {
	i32Pid := int32(pid)
	p, err := process.NewProcess(i32Pid)
	if err != nil {
		return false, err
	}
	// 进程名包含 name
	pName, err := p.Name()
	if err != nil {
		return false, err
	}
	fmt.Printf("process %d name %s\n", pid, pName)
	return strings.Contains(pName, name), nil
}
