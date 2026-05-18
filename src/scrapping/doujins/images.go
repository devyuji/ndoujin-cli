package doujins

import (
	"fmt"
	"net/http"
	"strings"

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
	req, err := http.NewRequest(http.MethodGet, c.Url, nil)

	if err != nil {
		return images, "", err
	}

	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	res, err := c.Client.Do(req)

	if err != nil {
		return images, "", err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return images, "", fmt.Errorf("Unable to access website %s - %d", c.Url, res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return images, "", err
	}

	folder := doc.Find(".folder-title").Children().Last()
	t := folder.Text()

	title := strings.TrimSpace(t)

	if title == "" {
		return images, "", fmt.Errorf("Title not found.")
	}

	doc.Find(".swiper-slide").Each(func(i int, s *goquery.Selection) {
		s.Find("img").Each(func(ii int, s *goquery.Selection) {
			src, found := s.Attr("data-src")

			if !found {
				return
			}

			u := strings.TrimSpace(src)

			images.Details = append(images.Details, types.ImagesDetails{
				Url: strings.TrimSpace(u),
			})
		})
	})

	return images, title, nil
}
