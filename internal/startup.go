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
	installed := InstallServiceIfNotYetInstalled(c)
	if installed {
		return true
	}

	register(c)
	removeExistingOldVersions(c)

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
