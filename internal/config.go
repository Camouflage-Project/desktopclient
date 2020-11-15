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
	SshServer SshServer `json:"ssh_server"`
	Forwards  []Forward `json:"forwards"`
}

type SshServer struct {
	Address  string `json:"address"`
	Username string `json:"username"`
}

type Forward struct {
	// local service to be forwarded
	Local Endpoint `json:"local"`
	// remote forwarding port (on remote SSH server network)
	Remote Endpoint `json:"remote"`
}

type Endpoint struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
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
		SshServer{
			Address:  "116.203.232.229:22",
			Username: "root",
		},
		[]Forward{
			{
				Local: Endpoint{Host: "127.0.0.1", Port: 8080},
				Remote: Endpoint{Host: "0.0.0.0", Port: 8119},
			},
		},
	}
}
