package internal

import (
	"desktopClient/config"
	"go.uber.org/zap"
	"strings"
	"time"
)

func InitHeartbeat(c *config.Configuration, logger *zap.Logger) {
	for {
		sendHeartbeat(c, logger)
		time.Sleep(5 * time.Second)
	}
}

func sendHeartbeat(c *config.Configuration, logger *zap.Logger) {
	ipParam := resolveIpParam(GetCurrentIp(logger))

	SendHeartbeatToBackend(c, ipParam, logger)
}

func resolveIpParam(ip string, err error) string {
	var ipParam string
	if err != nil {
		ipParam = err.Error()
	} else {
		ipParam = ip
	}

	return strings.TrimSpace(ipParam)
}
