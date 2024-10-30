package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ndoujin-cli",
	Short: "Download doujin from Nhentai website",
	Long:  "Download doujin from Nhentai website",
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
