package hitomi

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/types"
)

type Call struct {
	Url     string
	Headers map[string]string
}

func (c *Call) GetImages() (types.Image, error) {

	var images types.Image

	// 3665772
	// https://hitomi.la/doujinshi/sinful-lust-espa%C3%B1ol-3665772.html#1

	urlParse, err := url.Parse(c.Url)

	if err != nil {
		return images, fmt.Errorf("Unable to parse URL.")
	}

	path := urlParse.Path
	re := regexp.MustCompile(`\d+`)

	imageID := re.FindString(path)

	if imageID == "" {
		return images, fmt.Errorf("Image ID not found.")
	}

	url := fmt.Sprintf("https://hitomi.la/reader/%s.html#1", imageID)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return images, err
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil || res.StatusCode != http.StatusOK {
		return images, fmt.Errorf("Unable to access hitomi.la. Try to add Cookie on Config.json file or --cookie=<value>")
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return images, err
	}

	fmt.Println(doc.Text())

	return images, nil

}
