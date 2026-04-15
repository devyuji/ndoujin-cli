package nhentai

import (
	"fmt"
	"regexp"
)

func GetCode(url string) (string, error) {
	re := regexp.MustCompile(`/g/(\d+)`)
	match := re.FindStringSubmatch(url)
	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("unable.find.code")
}
