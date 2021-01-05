package util

import (
	"errors"
	"go.uber.org/zap"
	"net"
	"os/exec"
	"path/filepath"
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
	if err != nil {
		logger.Error(err.Error())
	}
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
