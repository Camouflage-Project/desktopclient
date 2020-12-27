package internal

import (
	"go.uber.org/zap"
	"os/exec"
	"time"
)

func InitScriptExecutor(c *Configuration, logger *zap.Logger) {
	for  {
		executeScript(c, logger)
		time.Sleep(60 * time.Second)
	}
}

func executeScript(c *Configuration, logger *zap.Logger) {
	script, err := FetchScriptFromBackend(c, logger)

	if err != nil {
		logger.Error(err.Error())
		return
	}
	if script == "" {
		return
	}

	output, err := exec.Command("/bin/bash", "-c", script).Output()

	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(string(output))
}

