package main

import (
	"desktopClient/config"
	"desktopClient/internal"
)

func main() {
	c := config.ReadConfig()

	stdLogger, logger := config.GetLoggers(c)

	logger.Info("starting...")

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
