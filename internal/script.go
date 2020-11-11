package internal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

func ExecuteScript(c *Configuration) {
	script, err := fetchScript(c)

	if err != nil {
		return
	}

	output, err := exec.Command("/bin/bash", "-c", script).Output()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))
}

func fetchScript(c *Configuration) (string, error) {
	resp, err := http.Get(c.ScriptUrl)
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
