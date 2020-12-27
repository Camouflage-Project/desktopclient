package internal

import (
	"errors"
	"go.uber.org/zap"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

func GetFilenameFromProcessName(processName string) string {
	return filepath.Base(processName)
}

func GetWorkingDirectoryFromProcessName(processName string) string {
	return filepath.Dir(processName) + string(filepath.Separator)
}

func ExecuteNewBinary(binaryPath string, logger *zap.Logger) {
	err := exec.Command(binaryPath).Start()
	if err!=nil {
		logger.Error(err.Error())
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

func GetOpenPort() (int, error) {
	for port := 9007; port <= 11000; port++ {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", strconv.Itoa(port)), timeout)
		if conn != nil {
			conn.Close()
		}
		if err != nil {
			return port, nil
		}
	}
	return 0, errors.New("no open port found")
}