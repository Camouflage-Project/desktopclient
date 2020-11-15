package main

import (
	"desktopClient/internal"
	"desktopClient/proxy"
	"desktopClient/tunnel"
)

func main() {
	c := internal.ReadConfig()
	installedService := internal.Startup(c)
	if installedService {
		return
	}

	done := make(chan bool)
	go internal.UpdateIfNewVersionExists(c)
	go internal.InitHeartbeat(c)

	proxy.InitializeForwardProxy()
	tunnel.InitializeTunnel(c)
	<- done
}


