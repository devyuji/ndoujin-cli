package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadImage(client *http.Client, url string, folderDir string, index int, headers map[string]string) {

	e := map[string]string{
		"image/jpeg": "jpg",
		"image/webp": "webp",
		"image/png":  "png",
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println("Invalid URL")
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("Unable to download %v\n", err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		fmt.Printf("Unable to access url: %s - %d\n", url, res.StatusCode)
		return
	}

	contentType := res.Header.Get("content-type")

	filename := fmt.Sprintf("%d.%s", index, e[contentType])

	filePath := filepath.Join(folderDir, filename)
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
