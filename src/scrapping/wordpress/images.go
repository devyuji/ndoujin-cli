package wordpress

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/types"
	"github.com/devyuji/ndoujin-cli/src/utils"
)

type Call struct {
	Client  *http.Client
	Url     string
	Headers map[string]string
}

func (c *Call) GetImages() (types.Image, string, error) {
	var images types.Image
	var f string

	req, err := http.NewRequest(http.MethodGet, c.Url, nil)

	if err != nil {
		return images, f, err
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	res, err := c.Client.Do(req)

	if err != nil {
		return images, f, err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return images, f, err
	}

	f = strings.TrimSpace(doc.Find("#chapter-heading").Text())

	doc.Find(".page-break").Each(func(i int, s *goquery.Selection) {
		v, exists := s.Find("img").Attr("src")

		if !exists {
			return
		}

		images.Details = append(images.Details, types.ImagesDetails{
			Url: strings.TrimSpace(v),
		})
	})

	f = utils.SanitizeFilename(f)

	return images, f, nil
}
