package internal

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

func GetFilenameFromProcessName(processName string) string {
	return filepath.Base(processName)
}

func GetWorkingDirectoryFromProcessName(processName string) string {
	return filepath.Dir(processName) + string(filepath.Separator)
}

func ExecuteNewBinary(binaryPath string) {
	err := exec.Command(binaryPath).Start()
	if err!=nil {
		fmt.Println(err.Error())
	}
}