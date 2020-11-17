package internal

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func GetInstallDirForOs(c *Configuration) string {
	var installDir string
	if runtime.GOOS == "windows" {
		installDir = c.WindowsInstallDirectory
	} else {
		installDir = c.UnixInstallDirectory
	}
	return installDir
}

func isSuperUser(c *Configuration, logger *zap.Logger) bool {
	installDir := GetInstallDirForOs(c)
	testFilePath := installDir + "desktopClientTestFile.txt"

	file, err := os.Create(testFilePath)
	if err != nil {
		logger.Error(err.Error())
		return false
	}

	err = file.Close()
	if err != nil {
		logger.Error(err.Error())
	}

	err = os.Remove(testFilePath)
	if err != nil {
		logger.Error(err.Error())
	}
	return true
}