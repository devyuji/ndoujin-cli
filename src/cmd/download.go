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
	"github.com/devyuji/ndoujin-cli/src/scrapping/myhentaigallery"
	"github.com/devyuji/ndoujin-cli/src/scrapping/nhentai"
	"github.com/devyuji/ndoujin-cli/src/scrapping/nhentaixxx"
	"github.com/devyuji/ndoujin-cli/src/scrapping/wordpress"
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

type provider interface {
	GetImages() (types.Image, string, error)
}

func init() {
	cmd.PersistentFlags().StringP("path", "p", "", "Set Download Path")
	cmd.PersistentFlags().BoolP("fail", "f", false, "Download failed ones.")

	rootCmd.AddCommand(cmd)
}

// n <- current t <- total
func printDownloadStatus(n int, t int) {
	fmt.Printf("\r\033[KDownloading image %d/%d...", n, t)
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func codeCmd(c *cobra.Command, args []string) {

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

				s := string(buffer[:bytesRead])
				for i := range strings.FieldsSeq(s) {
					start(c, i)
				}
			}

			return
		}

	}
	//-------------------END: look for code.txt file if no url is present -----------------------

	// URL Argument
	u := args[0]

	start(c, u)
}

func start(c *cobra.Command, uri string) {
	var hostName string
	var folderName string
	var headers = map[string]string{
		"User-Agent": config.Value.UserAgent,
	}

	// ------------------- Getting flags values -------------------------
	path, err := c.Flags().GetString("path")

	if err != nil {
		fmt.Println("Error: Getting flags...")
		return
	}

	fail, err := c.Flags().GetBool("fail")

	if err != nil {
		fmt.Println("Error: Getting flags...")
		return
	}
	//-------------------- END: Getting flags values ----------------------

	if path == "" {
		path = config.Value.Path
	}

	// Checking if url is valid or not
	if !isURL(uri) {
		fmt.Println("Please Enter Valid URL.")
		return
	}

	parseUrl, err := url.Parse(uri)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Fetching images for %s...\n", uri)

	var scrapper provider

	client := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: 60 * time.Second,
		},
		Timeout: 5 * time.Minute,
	}

	hostName = parseUrl.Host

	// ---------------------- Provider Config -------------------------
	switch hostName {
	case "nhentai.net":

		headers["Cookie"] = config.Value.Cookies.Nhentai

		scrapper = &nhentai.Call{
			Client:  client,
			Url:     uri,
			Headers: headers,
		}

	case "nhentai.xxx":

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

		scrapper = &doujins.Call{
			Client:  client,
			Url:     uri,
			Headers: headers,
		}

	case "www.myhentaigallery.com", "myhentaigallery.com":

		h := map[string]string{
			"Origin": "https://cdn.myhentaicomics.com",
			"Host":   "cdn.myhentaicomics.com",
			"Cookie": config.Value.Cookies.MyHentaiGallery,
		}

		maps.Copy(headers, h)

		scrapper = &myhentaigallery.Call{
			Client:  client,
			Url:     uri,
			Headers: headers,
		}

	case "hentaidoujinworld.com", "www.hentaidoujinworld.com":

		headers["Cookie"] = config.Value.Cookies.Hentaidoujinworld

		scrapper = &wordpress.Call{
			Client:  client,
			Url:     uri,
			Headers: headers,
		}

	default:
		fmt.Println("URL Not Supported.")
		return
	}

	// ------------- END: Provider Config ----------------------
	images, folderName, err := scrapper.GetImages()

	if err != nil {
		fmt.Println(err)
		return
	}

	if fail {
		err = readFailedDownload(client, folderName, headers, images)

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	totalImage := len(images.Details)

	if totalImage < 1 {
		fmt.Println("No images found. :-(")
		return
	}

	fmt.Printf("Total images found: %d\n", totalImage)

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
		countFailed := 0

		go func(index int) {
			defer wg.Done()
			defer func() { <-limiter }()

			printDownloadStatus(index, totalImage)
			err = utils.DownloadImage(client, detail.Url, downloadPath, index, headers)

			if err != nil {
				fmt.Printf("\n%v", err)
				countFailed++

				d := fmt.Sprintf("%d\n", index)
				saveFailedDownoad(d, folderName)
			}
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

func readFailedDownload(client *http.Client, name string, headers map[string]string, images types.Image) error {
	f := filepath.Join(config.Value.Path, fmt.Sprintf("%s.ndoujin", name))

	_, err := os.Stat(f)

	if os.IsNotExist(err) {
		return err
	}

	file, err := os.Open(f)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		file.Close()
		return err
	}

	scanner := bufio.NewScanner(file)
	var failedDownloads []string
	totalRetry := 10
	count := 1

	fmt.Println("Downloading failed images...")

	// scanner.Scan() returns false when it hits the end of the file
	for scanner.Scan() {
		retry := 1
		line := scanner.Text()

		if line == "" {
			continue
		}

		i, err := strconv.Atoi(line)

		if err != nil {
			fmt.Println(err)
			continue
		}

		url := images.Details[i-1]

		for retry < totalRetry {
			printDownloadStatus(count, -1)
			fmt.Println(" --Retry - ", retry)
			err := utils.DownloadImage(client, url.Url, filepath.Join(config.Value.Path, name), i, headers)

			if err != nil {
				fmt.Printf("\n%v", err)
				if retry == 1 {
					failedDownloads = append(failedDownloads, fmt.Sprintf("%d\n", i))
				}
				retry++
				continue
			}

			count++
			break
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
		file.Close()
		return err
	}

	file.Close()

	err = os.Remove(f)

	if err != nil {
		fmt.Println("Unable to remove file ", err)
	}

	for _, i := range failedDownloads {
		saveFailedDownoad(i, name)
	}

	return nil
}
