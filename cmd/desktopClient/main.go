package main

import (
	"desktopClient/internal"
	"desktopClient/proxy"
	"desktopClient/tunnel"
	"strconv"
)

func main() {
	c := internal.ReadConfig()

	stdLogger, logger := internal.GetLoggers(c)

	logger.Info("starting up desktopClient")
	logger.Info("injected key: " + c.Key)
	logger.Info("injected port: " + strconv.Itoa(c.Forwards[0].Remote.Port))

	installedService := internal.Install(c, logger)
	if installedService {
		return
	}

	err := internal.Register(c, logger)
	if err != nil {
		logger.Error("failed registration phase")
		return
	}

	done := make(chan bool)
	go internal.UpdateIfNewVersionExists(c, logger)
	go internal.InitHeartbeat(c, logger)
	go internal.InitScriptExecutor(c, logger)

	go proxy.InitializeForwardProxy(stdLogger, logger)
	go tunnel.InitializeTunnel(c, logger)
	logger.Info("everything initialized")
	<- done

	logger.Info("done and exiting")
}