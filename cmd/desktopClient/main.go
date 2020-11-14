package main

import (
	"desktopClient/internal"
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
	<- done
}


