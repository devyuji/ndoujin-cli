package doujins

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/devyuji/ndoujin-cli/src/types"
)

type Call struct {
	Body []byte
}

func (c *Call) GetImages() (types.Image, error) {
	var images types.Image

	b := bytes.NewReader(c.Body)

	doc, err := goquery.NewDocumentFromReader(b)

	if err != nil {
		return images, err
	}

	doc.Find(".swiper-slide").Each(func(i int, s *goquery.Selection) {
		s.Find("img").Each(func(ii int, s *goquery.Selection) {
			src, found := s.Attr("data-src")

			if !found {
				return
			}

			u := strings.TrimSpace(src)
			fmt.Println(u)

			images.Details = append(images.Details, types.ImagesDetails{
				Url:      strings.TrimSpace(u),
				FileName: fmt.Sprintf("%d.jpg", i),
			})
		})
		fmt.Println(i)
	})

	return images, nil
}
