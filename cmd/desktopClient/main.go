package main

import (
	"desktopClient/internal"
)

func main() {
	c := internal.ReadConfig()
	internal.Startup(c)

	done := make(chan bool)
	go internal.UpdateIfNewVersionExists(c)
	go internal.InitHeartbeat(c)
	<- done
}


