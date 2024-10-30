package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var USERAGENT = ""
var COOKIE = ""
var DOWNLOADPATH = ""

type Config struct {
	Path      string            `json:"path"`
	UserAgent string            `json:"user-agent"`
	Cookies   map[string]string `json:"cookies"`
}

func init() {
	file, err := os.Open("config.json")

	if err != nil {
		fmt.Println("no config.json file found!")
		return
	}

	data, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("error reading config.json file")
		return
	}

	var config Config

	err = json.Unmarshal(data, &config)

	if err != nil {
		fmt.Println("error reading config.json file")
		return
	}

	DOWNLOADPATH = config.Path
	USERAGENT = config.UserAgent

	for key, value := range config.Cookies {
		val := fmt.Sprintf("%s=%s; ", key, value)

		COOKIE += val
	}
}
