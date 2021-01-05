package internal

import (
	"desktopClient/config"
	"desktopClient/util"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func Install(c *config.Configuration, logger *zap.Logger, stdLogger *log.Logger) bool {
	if !isSuperUser(c, logger) {
		logSudoRequirementAndExit(logger)
	}

	copied := copyToInstallDirectoryAndExecute(c, logger)
	if copied {
		return true
	}

	installed := InstallServiceIfNotYetInstalled(c, logger)
	if installed {
		return true
	}

	return false
}

func Register(c *config.Configuration, logger *zap.Logger) error {
	err := RegisterOnBackend(c, logger)
	if err != nil {
		return err
	}

	err = removeExistingOldVersions(c, logger)
	if err != nil {
		return err
	}
	return nil
}

func copyToInstallDirectoryAndExecute(c *config.Configuration, logger *zap.Logger) bool {
	executable, err := os.Executable()
	if err != nil {
		logger.Error(err.Error())
	}

	fileName := util.GetFilenameFromProcessName(executable)
	workingDir := util.GetWorkingDirectoryFromProcessName(executable)

	installDir := getInstallDirForOs(c)

	newPath := installDir + fileName

	if workingDir != installDir {
		err = os.Rename(fileName, newPath)
		if err != nil {
			logger.Error(err.Error())
		}

		util.ExecuteNewBinary(newPath, logger)
		logger.Info("copied to " + newPath)
		return true
	}

	return false
}

func isSuperUser(c *config.Configuration, logger *zap.Logger) bool {
	installDir := getInstallDirForOs(c)
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

func getInstallDirForOs(c *config.Configuration) string {
	var installDir string
	if runtime.GOOS == "windows" {
		installDir = c.WindowsInstallDirectory
	} else {
		installDir = c.UnixInstallDirectory
	}
	return installDir
}

func removeExistingOldVersions(c *config.Configuration, logger *zap.Logger) error {
	processes, err := process.Processes()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			logger.Error(err.Error())
			return err
		} else if strings.HasPrefix(n, c.NamePrefix) && n != c.CurrentVersion {
			logger.Info("found something to kill!")
			err := removeFile(n, logger)
			if err != nil {
				return err
			}
			err = killExistingProcess(p, logger)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func removeFile(fileName string, logger *zap.Logger) error {
	err := os.Remove(fileName)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func killExistingProcess(p *process.Process, logger *zap.Logger) error {
	if runtime.GOOS == "windows" {
		err := killProcessOnWindows(int(p.Pid))
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	} else {
		err := p.Kill()
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	return nil
}

func killProcessOnWindows(p int) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(p))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}

func logSudoRequirementAndExit(logger *zap.Logger) {
	var errorMessage string
	if runtime.GOOS == "windows" {
		errorMessage = "Failed to start. Please run as administrator."
	} else {
		errorMessage = "Failed to start. Please execute with sudo."
	}

	logger.Error(errorMessage)
	os.Exit(1)
}
