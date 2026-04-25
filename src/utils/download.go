package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func DownloadImage(url string, folderDir string, fileName string, headers map[string]string) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println("Invalid URL")
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: time.Minute * 2,
		},
	}

	res, err := httpClient.Do(req)

	if err != nil {
		fmt.Printf("Unable to download %s\n", fileName)
		return
	}

	defer res.Body.Close()

	filePath := filepath.Join(folderDir, fileName)
	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}

	_, err = io.Copy(file, res.Body)

	if err != nil {
		fmt.Println("Failed to save image:", err)

		file.Close()
		return
	}

	if err = file.Sync(); err != nil {
		fmt.Println("Failed to save image:", err)

		file.Close()
		return
	}

	err = file.Close()

	if err != nil {
		fmt.Println("Failed to save image: ", err)
	}
}
