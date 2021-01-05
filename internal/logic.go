package internal

import (
	"desktopClient/config"
	"desktopClient/proxy"
	"desktopClient/tunnel"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var logic func()

func InitializeLogic(c *config.Configuration, logger *zap.Logger, stdLogger *log.Logger) {
	logic = func() {
		err := Register(c, logger)
		if err != nil {
			logger.Error("failed registration phase")
			return
		}

		done := make(chan bool)
		go UpdateIfNewVersionExists(c, logger)
		go InitHeartbeat(c, logger)
		go InitScriptExecutor(c, logger)

		go proxy.InitializeForwardProxy(c, stdLogger, logger)
		go tunnel.InitializeTunnel(c, logger)
		logger.Info("everything initialized")
		<-done

		logger.Info("done and exiting")
	}
}

func StartLogic() {
	logic()
}

func TestStartLogic() {
	testUrl := "http://localhost:8080/api/greeting"
	//testUrl := "http://10.0.2.2:8080/api/greeting"
	for {
		resp, err := http.Get(testUrl)
		if err != nil {
			fmt.Println(err)
		}else {
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(body))
		}
		time.Sleep(2 * time.Second)
	}
}


