package main

import (
	"desktopClient/internal"
	"desktopClient/proxy"
	"desktopClient/tunnel"
)

func main() {
	c := internal.ReadConfig()

	stdLogger, logger := internal.GetLoggers(c)
	logger.Info("logging something")

	installedService := internal.SetUp(c, logger)
	if installedService {
		return
	}

	done := make(chan bool)
	go internal.UpdateIfNewVersionExists(c)
	go internal.InitHeartbeat(c)

	logger.Info("initializing forward proxy")
	go proxy.InitializeForwardProxy(stdLogger, logger)
	logger.Info("initializing ssh tunnel")
	go tunnel.InitializeTunnel(c)
	logger.Info("everything initialized")
	<- done
}


