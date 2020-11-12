package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func InitHeartbeat(c *Configuration) {
	for  {

		sendHeartbeat(c)

		time.Sleep(10 * time.Second)
	}
}

func sendHeartbeat(c *Configuration) {
	ipParam := resolveIpParam(getCurrentIp())

	values := map[string]string{"key": c.Key, "ip": ipParam}

	jsonValue, _ := json.Marshal(values)

	resp, err := http.Post(c.NewVersionUrl,
		"application/json",
		bytes.NewBuffer(jsonValue))

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
	}
}

func getCurrentIp() (string, error) {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		fmt.Println(err)
		return "", err
	}else {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
			return "", err
		}

		return string(body), nil
	}
}

func resolveIpParam(ip string, err error) string {
	var ipParam string
	if err != nil {
		ipParam = err.Error()
	} else {
		ipParam = ip
	}

	return ipParam
}
