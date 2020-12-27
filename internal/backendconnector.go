package internal

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func RegisterOnBackend(c *Configuration, logger *zap.Logger) error {
	logger.Info("Registering with backend")
	values := map[string]string{"key": c.Key}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(c.RegistrationUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	defer resp.Body.Close()

	return nil
}

func GetNewestVersionFromBackend(c *Configuration, logger *zap.Logger) (string, error) {
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

func DownloadNewBinaryFromBackend(c *Configuration) (*http.Response, error) {
	values := map[string]string{"key": c.Key}
	jsonValue, _ := json.Marshal(values)

	return http.Post(c.BinaryUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))
}

func GetCurrentIp(logger *zap.Logger) (string, error) {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			logger.Error(err.Error())
			return "", err
		}

		return string(body), nil
	}
}

func SendHeartbeatToBackend(c *Configuration, ipParam string, logger *zap.Logger) {
	values := map[string]string{"key": c.Key, "ip": ipParam}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(c.HeartbeatUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	defer resp.Body.Close()

	if err != nil {
		logger.Error(err.Error())
	}
}

func FetchScriptFromBackend(c *Configuration, logger *zap.Logger) (string, error) {
	resp, err := http.Get(c.ScriptUrl)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			logger.Error(err.Error())
			return "", err
		}

		return string(body), nil
	}
}
