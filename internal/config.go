package internal

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"runtime"
	"strconv"
)

var Key = "MY8#m6P6hvQot%TJ1l7JLM"
var InjectedRemoteSshPort = "8119"

type Configuration struct {
	RegistrationUrl string
	ScriptUrl               string
	NewVersionUrl           string
	BinaryUrl               string
	HeartbeatUrl            string
	Key                     string
	NamePrefix              string
	CurrentVersion          string
	UnixInstallDirectory    string
	WindowsInstallDirectory string
	VerboseLogging          bool
	UnixLogOutputPath       string
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
	baseUrl := "http://localhost:8080/"
	executable, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	currentVersion := GetFilenameFromProcessName(executable)

	remoteSshPort, err := strconv.Atoi(InjectedRemoteSshPort)
	if err != nil {
		fmt.Println(err)
	}

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
		true,
		"/var/log/desktopClient.log",
		SshServer{
			Address:  "116.203.232.229:22",
			Username: "root",
			Password: "08Ih4rI$Nm6aNQx2%G*Drx*h9cR!gn5sk8kKEX^#rmf4cjM@6d",
		},
		[]Forward{
			{
				Local: Endpoint{Host: "127.0.0.1", Port: 8118},
				Remote: Endpoint{Host: "0.0.0.0", Port: remoteSshPort},
			},
		},
	}
}

func GetLoggers(config *Configuration) (*log.Logger, *zap.Logger) {
	var logPath string
	if runtime.GOOS == "windows" {
		logPath = config.WindowsInstallDirectory
	} else {
		logPath = config.UnixLogOutputPath
	}

	c := zap.NewProductionConfig()
	c.OutputPaths = []string{
		logPath,
	}
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if config.VerboseLogging {
		c.Level.SetLevel(zapcore.DebugLevel)
	} else {
		c.Level.SetLevel(zapcore.ErrorLevel)
	}

	logger, err := c.Build()
	if err != nil {
		log.Fatalln("Error: failed to initiate logger")
	}
	defer logger.Sync()
	return zap.NewStdLog(logger), logger
}
