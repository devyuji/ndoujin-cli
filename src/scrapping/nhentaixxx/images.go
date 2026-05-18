package nhentaixxx

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/config"
	"github.com/devyuji/ndoujin-cli/src/scrapping/nhentai"
	"github.com/devyuji/ndoujin-cli/src/types"
)

type Call struct {
	Client  *http.Client
	Url     string
	Headers map[string]string
}

func (c *Call) GetImages() (types.Image, string, error) {
	var images types.Image
	code, err := nhentai.GetCode(c.Url)

	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://nhentai.xxx/g/%s", code)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println(url)
		return images, "", fmt.Errorf("Invalid URL")
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	res, err := c.Client.Do(req)

	if err != nil {
		return images, "", err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return images, "", fmt.Errorf("Unable to access website : %s - %s", c.Url, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return images, "", fmt.Errorf("unable.to.parse.website")
	}

	totalPageStr := doc.Find(".pages").Text()
	totalPage, err := strconv.Atoi(totalPageStr)

	if err != nil {
		return images, "", err
	}

	limiter := make(chan int, config.Value.Concurrency)

	var wg sync.WaitGroup

	for i := range totalPage {
		limiter <- 1

		wg.Go(func() {
			image, err := c.getURL(code, i+1)

			if err != nil {
				fmt.Println("Unable to download image: ", err)

				<-limiter
			}

			images.Details = append(images.Details, image)

			<-limiter
		})

	}

	wg.Wait()

	return images, code, nil
}
