package internal

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UpdateIfNewVersionExists(c *Configuration, logger *zap.Logger) {
	logger.Info("initializing updater")
	for  {
		time.Sleep(2 * time.Second)

		newVersion, err := getNewestVersion(c, logger)
		if err != nil {
			continue
		}

		if newVersion != c.CurrentVersion {
			logger.Info("new version exists")
			filePath, err := downloadNewBinary(c, newVersion, logger)
			if err != nil {
				continue
			}
			ExecuteNewBinary(filePath)
		}
	}
}

func getNewestVersion(c *Configuration, logger *zap.Logger) (string, error) {
	values := map[string]string{"key": c.Key}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(c.NewVersionUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	if err != nil {
		logger.Error(err.Error())
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	return string(body), nil
}

func downloadNewBinary(c *Configuration, newVersion string, logger *zap.Logger) (string, error) {
	values := map[string]string{"key": c.Key, "binaryName": newVersion}
	jsonValue, _ := json.Marshal(values)

	response, err := http.Post(c.BinaryUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

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
