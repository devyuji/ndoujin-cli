package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/devyuji/ndoujin-cli/src/config"
)

func DownloadImage(url string, folderDir string, fileName string) {

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println("url is invalid")

		return
	}

	headers := map[string]string{
		"User-Agent": config.USERAGENT,
		"Cookies":    config.COOKIE,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{}

	res, err := httpClient.Do(req)

	if err != nil {
		fmt.Printf("unable to download %s\n", fileName)

		return
	}

	defer res.Body.Close()

	filePath := filepath.Join(folderDir, fileName)
	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println("failed to create file: %w", err)
		return
	}

	defer file.Close()

	_, err = io.Copy(file, res.Body)

	if err != nil {
		fmt.Println("failed to save image: %w", err)

		return
	}
}
