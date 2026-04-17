package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const VERSION string = "1.1.0"

type Cookie struct {
	Nhentai    string `json:"nhentai.net"`
	NhentaiXXX string `json:"nhentai.xxx"`
	Hitomi     string `json:"hitomi.la"`
}

type Config struct {
	Path      string `json:"path"`
	UserAgent string `json:"user-agent"`
	Cookies   Cookie `json:"cookies"`
}

var Value Config

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
