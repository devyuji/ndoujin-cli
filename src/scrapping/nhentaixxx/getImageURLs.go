package nhentaixxx

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/types"
)

func (c *Call) getURL(code string, pageNumber int) (types.ImagesDetails, error) {

	var imageDetails types.ImagesDetails

	url := fmt.Sprintf("https://nhentai.xxx/g/%s/%d", code, pageNumber)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return imageDetails, err
	}

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	res, err := c.Client.Do(req)

	if err != nil {
		return imageDetails, err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return imageDetails, fmt.Errorf("Unable to access website : %s - %d", c.Url, res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return imageDetails, err
	}

	imageUrl, exisit := doc.Find("#fimg").Attr("data-src")

	if !exisit {
		return imageDetails, fmt.Errorf("Image not found")
	}

	imageDetails.Url = imageUrl

	return imageDetails, nil
}
