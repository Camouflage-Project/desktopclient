package main

import (
	"desktopClient/config"
	"desktopClient/internal"
	"strconv"
)

func main() {
	c := config.ReadConfig()

	stdLogger, logger := config.GetLoggers(c)

	logger.Info("starting up desktopClient")
	logger.Info("injected key: " + c.ClientId)
	logger.Info("injected port: " + strconv.Itoa(c.Forwards[0].Remote.Port))
	logger.Info("injected NodeIp: " + c.SshServer.Address)
	logger.Info("injected NodeLimitedUserName: " + c.SshServer.Username)
	logger.Info("injected NodeLimitedUserPassword: " + c.SshServer.Password)

	internal.InitializeLogic(c, logger, stdLogger)
	if !c.RunAsBackgroundService {
		internal.StartLogic()
		return
	}

	installedService := internal.Install(c, logger, stdLogger)
	if installedService {
		return
	}

	internal.RunService(c, logger)
}
