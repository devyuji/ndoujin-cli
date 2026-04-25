package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const VERSION string = "v1.1.0"

type Cookie struct {
	Nhentai    string `json:"nhentai.net"`
	NhentaiXXX string `json:"nhentai.xxx"`
}

type Config struct {
	Path        string `json:"path"`
	UserAgent   string `json:"user-agent"`
	Cookies     Cookie `json:"cookies"`
	Concurrency int    `json:"concurrency"`
}

var Value *Config = &Config{
	Path:        "",
	UserAgent:   "",
	Cookies:     Cookie{},
	Concurrency: 10,
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

	err = json.Unmarshal(data, &Value)

	if err != nil {
		fmt.Println("error reading config.json file -", err)
		return
	}
}
