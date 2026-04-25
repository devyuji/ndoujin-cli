package cmd

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/devyuji/ndoujin-cli/src/config"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update ndoujin-cli to latest version",
	Long:    "Update ndoujin-cli to latest version",
	Example: "ndoujin-cli update",
	Run:     update,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

type githubApi struct {
	Name string `json:"name"`
	Url  string `json:"html_url"`
}

func update(c *cobra.Command, args []string) {
	c.Println("Checking for new version...")

	uri := "https://api.github.com/repos/devyuji/ndoujin-cli/releases/latest"

	res, err := http.Get(uri)

	if err != nil || res.StatusCode != http.StatusOK {
		c.Println("Unable to check for update: ", res.StatusCode)
		return
	}

	defer res.Body.Close()

	i, err := io.ReadAll(res.Body)

	if err != nil {
		c.Println("Unable to check for update: ", res.StatusCode)
		return
	}

	var response githubApi

	err = json.Unmarshal(i, &response)

	if err != nil {
		c.Println("Parse error")
		return
	}

	switch semver.Compare(config.VERSION, response.Name) {
	case -1:
		c.Println("Update is available download from here: ", response.Url)
	case 1, 0:
		c.Println("You're already up to date")
	}

}
