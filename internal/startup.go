package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func Startup(c *Configuration) bool {
	copied := copyToInstallDirectoryAndExecute(c)
	if copied {
		return true
	}

	installed := InstallServiceIfNotYetInstalled(c)
	if installed {
		return true
	}

	register(c)
	removeExistingOldVersions(c)

	return false
}

func copyToInstallDirectoryAndExecute(c *Configuration) bool {
	executable, err := os.Executable()
	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
		}

		ExecuteNewBinary(newPath)
		return true
	}

	return false
}

func removeExistingOldVersions(c *Configuration) {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println(err)
	}

	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			fmt.Println(err)
		} else if strings.HasPrefix(n, c.NamePrefix) && n != c.CurrentVersion {
			fmt.Println("found something to kill!")
			removeFile(n)
			killExistingProcess(p)
		}
	}
}

func removeFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		fmt.Println(err)
	}
}

func killExistingProcess(p *process.Process) {
	if runtime.GOOS == "windows" {
		err := killProcessOnWindows(int(p.Pid))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := p.Kill()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func register(c *Configuration) {
	values := map[string]string{"key": c.Key}

	jsonValue, _ := json.Marshal(values)

	_, err := http.Post(c.RegistrationUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Println(err)
	}
}

func killProcessOnWindows(p int) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(p))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}
