package nhentai

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/config"
)

type ImagesDetails struct {
	Url      string
	FileName string
}

type Image struct {
	Details []ImagesDetails
}

func GetImages(code string, saveInformation bool) (Image, error) {
	var images Image

	url := fmt.Sprintf("https://nhentai.net/g/%s", code)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		fmt.Println(url)
		return images, fmt.Errorf("url is invalid")
	}

	headers := map[string]string{
		"User-Agent": config.USERAGENT,
		"Cookies":    config.COOKIE,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)

	if err != nil || res.StatusCode != 200 {
		return images, fmt.Errorf("unablet to access website\nif the website is using cloudflare then add cookies in cookies.json file")
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return images, fmt.Errorf("unable to parse website")
	}

	var thumbImages []string

	doc.Find("#thumbnail-container").Find(".thumb-container").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Find("img").Attr("data-src")

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

		i := fmt.Sprintf("%s/%s/%s", imageBaseUrl, id, fileName)

		images.Details = append(images.Details, ImagesDetails{
			Url:      i,
			FileName: fileName,
		})
	}

	return images, nil
}
