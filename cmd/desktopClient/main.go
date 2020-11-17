package main

import (
	"desktopClient/internal"
	"desktopClient/proxy"
	"desktopClient/tunnel"
	"fmt"
)

func main() {
	c := internal.ReadConfig()

	fmt.Println(c.Key)
	fmt.Println(c.Forwards[0].Remote.Port)

	stdLogger, logger := internal.GetLoggers(c)

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


