package internal

import (
	"desktopClient/config"
	"desktopClient/util"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"time"
)

func UpdateIfNewVersionExists(c *config.Configuration, logger *zap.Logger) {
	logger.Info("initializing updater")
	for  {
		time.Sleep(10 * time.Second)

		newVersion, err := GetNewestVersionFromBackend(c, logger)
		if err != nil {
			continue
		}

		if newVersion != "" && newVersion != c.CurrentVersion {
			logger.Info("new version exists")
			filePath, err := downloadNewBinary(c, newVersion, logger)
			if err != nil {
				continue
			}
			util.ExecuteNewBinary(filePath, nil)
		}
	}
}

func downloadNewBinary(c *config.Configuration, newVersion string, logger *zap.Logger) (string, error) {
	response, err := DownloadNewBinaryFromBackend(c)

	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	defer response.Body.Close()

	file, err := os.Create(newVersion)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	err = file.Chmod(0700)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	binaryPath, err := filepath.Abs(file.Name())
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	return binaryPath, nil
}
