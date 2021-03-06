package config

import (
	"desktopClient/util"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"runtime"
	"strconv"
)

var ClientId = ""
var InjectedRemoteSshPort = ""
var BaseUrl = ""
var NodeIp = ""
var NodeLimitedUserName = ""
var NodeLimitedUserPassword = ""

type Configuration struct {
	RunAsBackgroundService  bool
	RegistrationUrl         string
	ScriptUrl               string
	NewVersionUrl           string
	BinaryUrl               string
	HeartbeatUrl            string
	ClientId                string
	NamePrefix              string
	CurrentVersion          string
	UnixInstallDirectory    string
	WindowsInstallDirectory string
	VerboseLogging          bool
	UnixLogOutputPath       string
	ProxyPort               int
	SshServer               SshServer `json:"ssh_server"`
	Forwards                []Forward `json:"forwards"`
}

type SshServer struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
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
	proxyPort := 10065
	executable, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	currentVersion := util.GetFilenameFromProcessName(executable)

	remoteSshPort, err := strconv.Atoi(InjectedRemoteSshPort)
	if err != nil {
		fmt.Println(err)
	}

	return &Configuration{
		true,
		BaseUrl + "/register-desktop-client",
		BaseUrl + "/script",
		BaseUrl + "/latest-version",
		BaseUrl + "/binary",
		BaseUrl + "/heartbeat",
		ClientId,
		"ResidentialProxyClient",
		currentVersion,
		"/usr/local/bin/",
		"%userprofile%",
		true,
		"/var/log/desktopClient.log",
		proxyPort,
		SshServer{
			Address:  NodeIp + ":22",
			Username: NodeLimitedUserName,
			Password: NodeLimitedUserPassword,
		},
		[]Forward{
			{
				Local:  Endpoint{Host: "127.0.0.1", Port: proxyPort},
				Remote: Endpoint{Host: "0.0.0.0", Port: remoteSshPort},
			},
		},
	}
}

func GetLoggers(config *Configuration) (*log.Logger, *zap.Logger) {
	var logPath string
	if runtime.GOOS == "windows" {
		logPath = config.WindowsInstallDirectory + "desktopClient.log"
	} else {
		logPath = config.UnixLogOutputPath
	}

	c := zap.NewProductionConfig()
	if runtime.GOOS == "windows" {
		c.OutputPaths = []string{
			"stdout",
		}
	} else {
		c.OutputPaths = []string{
			"stdout", logPath,
		}
	}

	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if config.VerboseLogging {
		c.Level.SetLevel(zapcore.DebugLevel)
	} else {
		c.Level.SetLevel(zapcore.ErrorLevel)
	}

	logger, err := c.Build()
	if err != nil {
		fmt.Println(err)
	}
	defer logger.Sync()
	return zap.NewStdLog(logger), logger
}
