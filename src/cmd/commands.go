package cmd

import (
	"fmt"
	"log"

	"github.com/devyuji/ndoujin-cli/src/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ndoujin-cli",
	Short: "Download doujin from Supported Websites",
	Long:  "Download doujin from Supported Websites",
	Run:   root,
}

func Execute() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number and exit")

	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}

func root(c *cobra.Command, args []string) {
	version, err := c.Flags().GetBool("version")

	if err != nil {
		c.Println("Unable to get flag version")
		return
	}

	if version {
		fmt.Println(config.VERSION)
	}

}
