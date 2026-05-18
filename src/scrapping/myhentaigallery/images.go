package myhentaigallery

import (
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/types"
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

	doc.Find(".comic-thumb").Each(func(i int, s *goquery.Selection) {
		image, e := s.Find("img").Attr("src")

		if !e {
			return
		}

		re := regexp.MustCompile(`/thumbnail`)

		newURL := re.ReplaceAllString(image, "/original")

		images.Details = append(images.Details, types.ImagesDetails{Url: newURL})
	})

	return images, f, nil

}
