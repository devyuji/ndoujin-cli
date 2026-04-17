package cmd

import (
	"fmt"
	"log"
	"maps"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/devyuji/ndoujin-cli/src/config"
	"github.com/devyuji/ndoujin-cli/src/scrapping/hitomi"
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
	Example: "ndoujin-cli download 533999\nndoujin-cli download https://nhentai.net/g/533999",
	Run:     codeCmd,
}

func init() {
	rootCmd.AddCommand(cmd)

	rootCmd.PersistentFlags().StringP("path", "p", "", "Set Download Path")
	rootCmd.PersistentFlags().StringP("cookie", "c", "", "Add Cookie")
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
		fmt.Println("Something went wrong!")
		return
	}

	// cookie, err := c.Flags().GetString("cookie")

	// if err != nil {
	// 	fmt.Println("Something went wrong!")
	// 	return
	// }

	if config.Value.Path != "" {
		path = config.Value.Path
	}

	// look for code.txt file if no code or url is present
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

	switch hostName {
	case "nhentai.net":
		folderName, err = nhentai.GetCode(uri)

		if err != nil {
			log.Fatal(err)
		}

		scrapper = &nhentai.Call{
			Url: uri,
		}

		headers["Cookie"] = config.Value.Cookies.Nhentai

	case "nhentai.xxx":
		folderName, err = nhentai.GetCode(uri)

		if err != nil {
			log.Fatal(err)
		}

		scrapper = &nhentaixxx.Call{
			Url: uri,
		}

		headers["Cookie"] = config.Value.Cookies.NhentaiXXX

	case "hitomi.la":
		fmt.Println("Under Construction!")
		folderName = "demo"

		h := map[string]string{
			"Referer": "https://hitomi.la/",
			"Cookie":  config.Value.Cookies.Hitomi,
		}

		maps.Copy(headers, h)

		scrapper = &hitomi.Call{
			Url:     uri,
			Headers: headers,
		}

	default:
		fmt.Println("URL Not Supported.")
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

	// ------------------ create folder for download ------------------------
	downloadPath := filepath.Join(path, folderName)

	_, err = os.Stat(folderName)

	if os.IsNotExist(err) {
		err = os.Mkdir(downloadPath, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

	// ------------------ downloading start here -------------------------
	limiter := make(chan int, 10)
	var wg sync.WaitGroup

	fmt.Printf("Downloading %s\n", uri)
	for _, detail := range images.Details {
		limiter <- 1

		wg.Go(func() {
			utils.DownloadImage(detail.Url, downloadPath, detail.FileName, headers)

			<-limiter
		})
	}

	wg.Wait()

	fmt.Println("Download completed.")
}
