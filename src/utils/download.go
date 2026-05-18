package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func getFileType(rawURL string) string {
	// 1. Parse the URL to ignore the domain and query parameters (?word=hello)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Invalid URL provided")
		return ""
	}

	// parsedURL.Path is now "/s/fgs/1.jpeg/signature"

	// 2. Split the path into individual segments
	segments := strings.Split(parsedURL.Path, "/")

	var fileType string

	// 3. Loop through segments to find the last one with a file extension
	for _, segment := range segments {
		// If the segment has a dot, assume it's a file
		if strings.Contains(segment, ".") {
			// Split by the dot and grab the last part
			parts := strings.Split(segment, ".")
			fileType = parts[len(parts)-1]
		}
	}

	return fileType
}

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
	extension := e[contentType]

	if contentType == "" {
		extension = getFileType(url)
	}

	filename := fmt.Sprintf("%d.%s", index, extension)

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
