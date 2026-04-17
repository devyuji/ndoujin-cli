package nhentaixxx

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/config"
	"github.com/devyuji/ndoujin-cli/src/scrapping/nhentai"
	"github.com/devyuji/ndoujin-cli/src/types"
)

type Call struct {
	Url string
}

func (c *Call) GetImages() (types.Image, error) {
	var images types.Image
	code, err := nhentai.GetCode(c.Url)

	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://nhentai.xxx/g/%s", code)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println(url)
		return images, fmt.Errorf("url is invalid")
	}

	headers := map[string]string{
		"User-Agent": config.Value.UserAgent,
		"cookie":     config.Value.Cookies.NhentaiXXX,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)

	if err != nil || res.StatusCode != 200 {
		return images, fmt.Errorf("unable to access website\nif the website is using cloudflare then add cookies in config.json file - %d", res.StatusCode)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return images, fmt.Errorf("unable.to.parse.website")
	}

	totalPageStr := doc.Find(".pages").Text()
	totalPage, err := strconv.Atoi(totalPageStr)

	if err != nil {
		return images, err
	}

	limiter := make(chan int, 10)

	var wg sync.WaitGroup

	for i := range totalPage {
		limiter <- 1

		wg.Go(func() {
			image, err := getUrl(code, i+1)

			if err != nil {
				fmt.Println("Unable to download image:", image.FileName)

				<-limiter
				fmt.Println("not working")
			}

			images.Details = append(images.Details, image)

			<-limiter
		})

	}

	wg.Wait()

	return images, nil
}

func getUrl(code string, pageNumber int) (types.ImagesDetails, error) {
	var imageDetails types.ImagesDetails

	url := fmt.Sprintf("https://nhentai.xxx/g/%s/%d", code, pageNumber)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return imageDetails, err
	}

	headers := map[string]string{
		"User-Agent": config.Value.UserAgent,
		"Cookie":     config.Value.Cookies.NhentaiXXX,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil || res.StatusCode != 200 {
		return imageDetails, fmt.Errorf("unable to access website\nif the website is using cloudflare then add cookies in config.json file - %d", res.StatusCode)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return imageDetails, err
	}

	imageUrl, exisit := doc.Find("#fimg").Attr("data-src")

	if !exisit {
		return imageDetails, fmt.Errorf("image.not.found")
	}

	imageDetails.FileName = getFileName(imageUrl)
	imageDetails.Url = imageUrl

	return imageDetails, nil
}

func getFileName(url string) string {
	urlSplit := strings.Split(url, "/")
	length := len(urlSplit)
	fileName := urlSplit[length-1]

	return fixFile(fileName)
}

func fixFile(input string) string {
	parts := strings.Split(input, ".")

	if len(parts) < 3 {
		return input
	}

	filename := strings.Join(parts[:len(parts)-1], ".")

	return filename
}
