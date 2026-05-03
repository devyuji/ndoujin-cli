package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadImage(client *http.Client, url string, folderDir string, index int, headers map[string]string) error {

	e := map[string]string{
		"image/jpeg": "jpg",
		"image/webp": "webp",
		"image/png":  "png",
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return fmt.Errorf("Invalid URL\n")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Unable to download %v\n", err)
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("Unable to access url: %s - %d\n", url, res.StatusCode)
	}

	contentType := res.Header.Get("content-type")

	filename := fmt.Sprintf("%d.%s", index, e[contentType])

	filePath := filepath.Join(folderDir, filename)
	file, err := os.Create(filePath)

	if err != nil {
		return fmt.Errorf("Failed to create file: %v\n", err)

	}

	_, err = io.Copy(file, res.Body)

	if err != nil {
		file.Close()
		return fmt.Errorf("Failed to save image: %v\n", err)
	}

	if err = file.Sync(); err != nil {
		file.Close()
		return fmt.Errorf("Failed to save image: %v\n", err)

	}

	err = file.Close()

	if err != nil {
		return fmt.Errorf("Failed to save image: %v\n", err)

	}

	return nil
}
