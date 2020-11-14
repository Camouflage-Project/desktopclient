package internal

import (
	"fmt"
	"os"
)

var Key = "MY8#m6P6hvQot%TJ1l7JLM"

type Configuration struct {
	RegistrationUrl string
	ScriptUrl       string
	NewVersionUrl   string
	BinaryUrl string
	HeartbeatUrl string
	Key string
	NamePrefix string
	CurrentVersion string
	UnixInstallDirectory string
	WindowsInstallDirectory string
}

func ReadConfig() *Configuration {
	baseUrl := "http://localhost:8080/"
	executable, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	currentVersion := GetFilenameFromProcessName(executable)

	return &Configuration{
		baseUrl + "register-desktop-client",
		baseUrl + "script",
		baseUrl + "new-version",
		baseUrl + "binary",
		baseUrl + "heartbeat",
		Key,
		"SingleProxyDesktopClient",
		currentVersion,
		"/usr/local/bin/",
		"C:\\Tools\\",
	}
}
