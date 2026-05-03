package cmd

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/devyuji/ndoujin-cli/src/config"
	"github.com/devyuji/ndoujin-cli/src/scrapping/doujins"
	"github.com/devyuji/ndoujin-cli/src/scrapping/nhentai"
	"github.com/devyuji/ndoujin-cli/src/scrapping/nhentaixxx"
	"github.com/devyuji/ndoujin-cli/src/types"
	"github.com/devyuji/ndoujin-cli/src/utils"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "download",
	Short:   "Download doujin",
	Long:    "Download doujin",
	Example: "ndoujin-cli download https://nhentai.net/g/533999",
	Run:     codeCmd,
}

func init() {
	cmd.PersistentFlags().StringP("path", "p", "", "Set Download Path")

	rootCmd.AddCommand(cmd)
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

type provider interface {
	GetImages() (types.Image, error)
}

func codeCmd(c *cobra.Command, args []string) {

	path, err := c.Flags().GetString("path")

	if err != nil {
		fmt.Println("Unable to get path flag")
		return
	}

	if path == "" {
		path = config.Value.Path
	}

	//---------------------- look for code.txt file if no url is present -----------------------
	if len(args) < 1 {
		_, err := os.Stat("code.txt")

		if !os.IsNotExist(err) {

			file, err := os.Open("code.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			buffer := make([]byte, 1024) // 1KB buffer size
			for {
				bytesRead, err := file.Read(buffer)
				if err != nil {
					if err.Error() == "EOF" { // End of file
						break
					}
					log.Fatal(err)
				}

				c := string(buffer[:bytesRead])
				for i := range strings.FieldsSeq(c) {
					start(i, path)
				}
			}

			return
		}

	}
	//---------------------- look for code.txt file if no url is present -----------------------

	u := args[0]

	start(u, path)
}

func start(uri string, path string) {
	var hostName string
	var folderName string
	var headers = map[string]string{
		"User-Agent": config.Value.UserAgent,
	}

	if !isURL(uri) {
		fmt.Println("Please Enter Valid URL.")
		return
	}

	parseUrl, err := url.Parse(uri)

	if err != nil {
		log.Fatal(err)
	}

	hostName = parseUrl.Host

	fmt.Printf("Fetching images for %s...\n", uri)

	var scrapper provider

	client := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 60 * time.Second,
		},
		Timeout: 5 * time.Minute,
	}

	switch hostName {
	case "nhentai.net":
		folderName, err = nhentai.GetCode(uri)

		if err != nil {
			fmt.Println(err)
			return
		}

		headers["Cookie"] = config.Value.Cookies.Nhentai

		scrapper = &nhentai.Call{
			Client:  client,
			Url:     uri,
			Headers: headers,
		}

	case "nhentai.xxx":
		folderName, err = nhentai.GetCode(uri)

		if err != nil {
			fmt.Println(err)
			return
		}

		headers["Cookie"] = config.Value.Cookies.NhentaiXXX

		scrapper = &nhentaixxx.Call{
			Client:  client,
			Url:     uri,
			Headers: headers,
		}

	case "doujins.com":
		h := map[string]string{
			"Host":            "static.doujins.com",
			"Accept-Language": "en-US,en;q=0.9",
			"Sec-GPC":         "1",
			"Connection":      "keep-alive",
			"Referer":         "https://doujins.com/",
			"Sec-Fetch-Dest":  "image",
			"Sec-Fetch-Mode":  "no-cors",
			"Sec-Fetch-Site":  "same-site",
			"DNT":             "1",
			"Priority":        "u=5, i",
			"Pragma":          "no-cache",
			"Cache-Control":   "no-cache",
			"Cookie":          config.Value.Cookies.Doujins,
		}

		maps.Copy(headers, h)

		d, b, err := doujins.GetDetails(client, uri, headers)

		if err != nil {
			fmt.Println("Unable to get details - ", err)
			return
		}

		folderName = d["name"].(string)

		scrapper = &doujins.Call{
			Body: b,
		}

	default:
		fmt.Println("URL Not Supported.")
		return
	}

	c := readFailedDownload(client, folderName, headers)

	if c {
		fmt.Println("Download Complete.")
		return
	}

	images, err := scrapper.GetImages()

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(images.Details) < 1 {
		fmt.Println("No images found. :-(")
		return
	}

	fmt.Printf("Total images found: %d\n", len(images.Details))

	// ------------------ creating folder  ------------------------
	downloadPath := filepath.Join(path, folderName)

	_, err = os.Stat(downloadPath)

	if os.IsNotExist(err) {
		err = os.Mkdir(downloadPath, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}
	// ------------------ creating folder ------------------------

	// ------------------ downloading images -------------------------
	limiter := make(chan int, config.Value.Concurrency)
	var wg sync.WaitGroup

	for i, detail := range images.Details {
		limiter <- 1
		wg.Add(1)

		go func(index int) {

			fmt.Printf("\r\033[KDownloading image %d/%d...", index, len(images.Details))
			err = utils.DownloadImage(client, detail.Url, downloadPath, index, headers)

			if err != nil {
				fmt.Println(err)

				d := fmt.Sprintf("%d;%s\n", index, detail.Url)
				saveFailedDownoad(d, folderName)
			}

			defer wg.Done()
			defer func() { <-limiter }()

		}(i + 1)
	}

	wg.Wait()

	fmt.Println("\nDownload completed.")
}

func saveFailedDownoad(data string, name string) {

	fileName := fmt.Sprintf("%s.ndoujin", name)

	f := filepath.Join(config.Value.Path, fileName)

	file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 4. Write the formatted string
	_, err = file.WriteString(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func readFailedDownload(client *http.Client, name string, headers map[string]string) bool {
	f := filepath.Join(config.Value.Path, fmt.Sprintf("%s.ndoujin", name))

	_, err := os.Stat(f)

	if os.IsNotExist(err) {
		return false
	}

	fmt.Println("Downloading Failed Images...")

	file, err := os.Open(f)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		file.Close()
		return false
	}

	scanner := bufio.NewScanner(file)
	var failedDownloads []string

	// scanner.Scan() returns false when it hits the end of the file
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		s := strings.Split(line, ";")
		fn, _ := strconv.Atoi(s[0])
		u := s[1]

		err := utils.DownloadImage(client, strings.TrimSpace(u), filepath.Join(config.Value.Path, name), fn, headers)

		if err != nil {
			fmt.Println(err)
			failedDownloads = append(failedDownloads, fmt.Sprintf("%d;%s", fn, u))
			file.Close()
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
		file.Close()
		return false
	}

	file.Close()

	if len(failedDownloads) < 1 {
		err := os.Remove(f)

		if err != nil {
			fmt.Println("Unable to remove file ", err)
		}
	}

	for _, i := range failedDownloads {
		saveFailedDownoad(i, name)
	}

	return true
}
