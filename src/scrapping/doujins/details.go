package doujins

import (
	"fmt"
	"net/http"
)

func (c *Call) getDetails(client *http.Client) (map[string]string, error) {
	fmt.Println(c)
	var details = make(map[string]string)

	// req, err := http.NewRequest(http.MethodGet, c.Url, nil)

	// if err != nil {
	// 	return details, err
	// }

	return details, nil
}
