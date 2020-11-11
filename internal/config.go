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
	Key string
	NamePrefix string
	CurrentVersion string
}

func ReadConfig() *Configuration {
	baseUrl := "http://localhost:8080/"
	executable, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	currentVersion := getFilenameFromProcessName(executable)

	return &Configuration{
		baseUrl + "register-desktop-client",
		baseUrl + "script",
		baseUrl + "new-version",
		baseUrl + "binary",
		Key,
		"SingleProxyDesktopClient",
		currentVersion,
	}
}
