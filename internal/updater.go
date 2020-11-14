package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UpdateIfNewVersionExists(c *Configuration) {
	for  {
		time.Sleep(2 * time.Second)

		newVersion, err := getNewestVersion(c)
		if err != nil {
			continue
		}

		if newVersion != c.CurrentVersion {
			fmt.Println("new version exists")
			filePath, err := downloadNewBinary(c, newVersion)
			if err != nil {
				continue
			}
			ExecuteNewBinary(filePath)
		}
	}
}

func getNewestVersion(c *Configuration) (string, error) {
	values := map[string]string{"key": c.Key}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(c.NewVersionUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

func downloadNewBinary(c *Configuration, newVersion string) (string, error) {
	values := map[string]string{"key": c.Key, "binaryName": newVersion}
	jsonValue, _ := json.Marshal(values)

	response, err := http.Post(c.BinaryUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer response.Body.Close()

	file, err := os.Create(newVersion)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	err = file.Chmod(0700)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	binaryPath, err := filepath.Abs(file.Name())
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return binaryPath, nil
}
