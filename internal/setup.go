package internal

import (
	"bytes"
	"encoding/json"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func SetUp(c *Configuration, logger *zap.Logger) bool {
	copied := copyToInstallDirectoryAndExecute(c, logger)
	if copied {
		return true
	}

	installed := InstallServiceIfNotYetInstalled(c)
	if installed {
		return true
	}

	register(c, logger)
	removeExistingOldVersions(c, logger)

	logger.Info("prerequisites completed")

	return false
}

func copyToInstallDirectoryAndExecute(c *Configuration, logger *zap.Logger) bool {
	executable, err := os.Executable()
	if err != nil {
		logger.Error(err.Error())
	}

	fileName := GetFilenameFromProcessName(executable)
	workingDir := GetWorkingDirectoryFromProcessName(executable)

	var installDir string
	if runtime.GOOS == "windows" {
		installDir = c.WindowsInstallDirectory
	} else {
		installDir = c.UnixInstallDirectory
	}

	newPath := installDir + fileName

	if workingDir != installDir {
		err = os.Rename(fileName, newPath)
		if err != nil {
			logger.Error(err.Error())
		}

		ExecuteNewBinary(newPath)
		return true
	}

	return false
}

func removeExistingOldVersions(c *Configuration, logger *zap.Logger) {
	processes, err := process.Processes()
	if err != nil {
		logger.Error(err.Error())
	}

	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			logger.Error(err.Error())
		} else if strings.HasPrefix(n, c.NamePrefix) && n != c.CurrentVersion {
			logger.Info("found something to kill!")
			removeFile(n, logger)
			killExistingProcess(p, logger)
		}
	}
}

func removeFile(fileName string, logger *zap.Logger) {
	err := os.Remove(fileName)
	if err != nil {
		logger.Error(err.Error())
	}
}

func killExistingProcess(p *process.Process, logger *zap.Logger) {
	if runtime.GOOS == "windows" {
		err := killProcessOnWindows(int(p.Pid))
		if err != nil {
			logger.Error(err.Error())
		}
	} else {
		err := p.Kill()
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func register(c *Configuration, logger *zap.Logger) {
	values := map[string]string{"key": c.Key}

	jsonValue, _ := json.Marshal(values)

	_, err := http.Post(c.RegistrationUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	if err != nil {
		logger.Error(err.Error())
	}
}

func killProcessOnWindows(p int) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(p))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}
