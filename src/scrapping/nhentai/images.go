package nhentai

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/config"
	"github.com/devyuji/ndoujin-cli/src/types"
)

type Call struct {
	Url string
}

func (c Call) Get() (types.Image, error) {
	var images types.Image
	code, err := GetCode(c.Url)

	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("https://nhentai.net/g/%s", code)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println(url)
		return images, fmt.Errorf("url is invalid")
	}

	headers := map[string]string{
		"User-Agent": config.USERAGENT,
		"cookie":     config.COOKIE,
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
		return images, fmt.Errorf("unable to parse website")
	}

	var thumbImages []string

	doc.Find("#thumbnail-container").Find(".thumb-container").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Find("img").Attr("src")

		if exists {
			thumbImages = append(thumbImages, src)
		}
	})

	imageBaseUrl := "https://i3.nhentai.net/galleries"

	for _, val := range thumbImages {
		urlSplit := strings.Split(val, "/")
		length := len(urlSplit)
		id := urlSplit[length-2]
		fileName := urlSplit[length-1]
		fileName = strings.Replace(fileName, "t", "", 1)

		i := fmt.Sprintf("%s/%s/%s", imageBaseUrl, id, fixFile(fileName))

		images.Details = append(images.Details, types.ImagesDetails{
			Url:      i,
			FileName: fixFile(fileName),
		})
	}

	return images, nil
}

func fixFile(input string) string {
	parts := strings.Split(input, ".")

	if len(parts) < 3 {
		return input
	}

	filename := strings.Join(parts[:len(parts)-1], ".")

	return filename
}
