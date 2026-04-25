package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ndoujin-cli",
	Short: "Download doujin from Supported Websites",
	Long:  "Download doujin from Supported Websites",
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}

}
