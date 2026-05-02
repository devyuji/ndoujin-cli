package doujins

import (
	"net/http"

	"github.com/devyuji/ndoujin-cli/src/types"
)

type Call struct {
	Url     string
	Headers map[string]string
}

func (c *Call) GetImages() (types.Image, error) {
	var images types.Image

	httpClient := &http.Client{}

	c.getDetails(httpClient)

	return images, nil
}
