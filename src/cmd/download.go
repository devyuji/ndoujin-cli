package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/devyuji/ndoujin-cli/src/config"
	"github.com/devyuji/ndoujin-cli/src/scrapping/nhentai"
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

func codeCmd(c *cobra.Command, args []string) {

	path, err := c.Flags().GetString("path")

	if err != nil {
		fmt.Println("Wrong flags")
		return
	}

	if config.DOWNLOADPATH != "" {
		path = config.DOWNLOADPATH
	}

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
				for _, i := range strings.Fields(c) {
					start(i, path)
				}
			}

			return
		}

	}

	u := args[0]

	var code string = u

	if isURL(u) {
		re := regexp.MustCompile(`/g/(\d+)`)
		match := re.FindStringSubmatch(u)
		if len(match) > 1 {
			code = match[1]
		}
	}

	start(code, path)
}

func start(code string, path string) {

	var images nhentai.Image

	images, err := nhentai.GetImages(code, false)

	if err != nil {
		log.Fatal(err)
	}

	downloadPath := filepath.Join(path, code)

	_, err = os.Stat(code)

	if os.IsNotExist(err) {
		err = os.Mkdir(downloadPath, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

	limiter := make(chan int, 10)
	var wg sync.WaitGroup

	fmt.Printf("Downloading %s\n", code)
	for _, detail := range images.Details {
		limiter <- 1

		wg.Add(1)

		go func() {
			defer wg.Done()

			utils.DownloadImage(detail.Url, downloadPath, detail.FileName)

			<-limiter
		}()
	}

	wg.Wait()

	fmt.Println("Download completed")
}
