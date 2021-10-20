package internal

import (
	"bytes"
	"desktopClient/config"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func RegisterOnBackend(c *config.Configuration, logger *zap.Logger) error {
	logger.Info("Registering with backend")
	values := map[string]string{"clientId": c.ClientId}

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

func GetNewestVersionFromBackend(c *config.Configuration, logger *zap.Logger) (string, error) {
	values := map[string]string{"clientId": c.ClientId}

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

func DownloadNewBinaryFromBackend(c *config.Configuration) (*http.Response, error) {
	values := map[string]string{"clientId": c.ClientId}
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

func SendHeartbeatToBackend(c *config.Configuration, ipParam string, logger *zap.Logger) {
	values := map[string]string{"clientId": c.ClientId, "ip": ipParam}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(c.HeartbeatUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer resp.Body.Close()
}

func FetchScriptFromBackend(c *config.Configuration, logger *zap.Logger) (string, error) {
	resp, err := http.Get(c.ScriptUrl)
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
