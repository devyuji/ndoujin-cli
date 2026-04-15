package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/devyuji/ndoujin-cli/src/config"
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

	rootCmd.PersistentFlags().StringP("path", "p", "", "enter download path")
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

type img interface {
	Get() (types.Image, error)
}

func codeCmd(c *cobra.Command, args []string) {

	path, err := c.Flags().GetString("path")

	if err != nil {
		fmt.Println("Wrong flags")
		return
	}

	if config.DOWNLOADPATH != "" {
		path = config.DOWNLOADPATH
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
	var code string
	var images types.Image
	var err error

	if !isURL(uri) {
		fmt.Println("invalid url")
		return
	}

	parseUrl, err := url.Parse(uri)

	if err != nil {
		log.Fatal(err)
	}

	hostName = parseUrl.Host

	fmt.Printf("Fetching images for %s...\n", uri)

	switch hostName {
	case "nhentai.net":
		code, err = nhentai.GetCode(uri)

		if err != nil {
			log.Fatal(err)
		}

		c := nhentai.Call{
			Url: uri,
		}

		images, err = img.Get(c)

	case "nhentai.xxx":
		code, err = nhentai.GetCode(uri)

		if err != nil {
			log.Fatal(err)
		}

		c := nhentaixxx.Call{
			Url: uri,
		}

		images, err = img.Get(c)

	default:
		log.Fatal("url.not.supported")
	}

	if len(images.Details) < 1 {
		fmt.Println("No images found. :-(")
		return
	}

	fmt.Printf("Total images found: %d\n", len(images.Details))

	// ------------------ create folder for download ------------------------
	downloadPath := filepath.Join(path, code)

	_, err = os.Stat(code)

	if os.IsNotExist(err) {
		err = os.Mkdir(downloadPath, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

	// ------------------ downloading start here -------------------------
	limiter := make(chan int, 10)
	var wg sync.WaitGroup

	fmt.Printf("Downloading %s\n", code)
	for _, detail := range images.Details {
		limiter <- 1

		wg.Go(func() {
			utils.DownloadImage(detail.Url, downloadPath, detail.FileName)

			<-limiter
		})
	}

	wg.Wait()

	fmt.Println("Download completed.")
}
