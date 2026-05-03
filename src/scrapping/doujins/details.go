package doujins

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetDetails(client *http.Client, uri string, headers map[string]string) (map[string]any, []byte, error) {
	var details = make(map[string]any)

	req, err := http.NewRequest(http.MethodGet, uri, nil)

	if err != nil {
		return details, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)

	if err != nil {
		return details, nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return details, nil, fmt.Errorf("Unable to access website %s - %d", uri, res.StatusCode)
	}

	i, err := io.ReadAll(res.Body)

	if err != nil {
		return details, nil, err
	}

	b := bytes.NewReader(i)

	d, err := goquery.NewDocumentFromReader(b)

	if err != nil {
		return details, nil, err
	}

	folder := d.Find(".folder-title").Children().Last()
	t := folder.Text()

	title := strings.TrimSpace(t)

	if title == "" {
		return details, nil, fmt.Errorf("Title not found.")
	}

	details["name"] = title

	return details, i, nil
}
