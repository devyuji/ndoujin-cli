package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ndoujin-cli",
	Short: "Download doujin from Supported Websites",
	Long:  "Download doujin from Supported Websites",
	Run:   root,
}

func Execute() {
	rootCmd.PersistentFlags().BoolP("update", "u", false, "Update Ndoujin-CLI to latest version.")

	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}

}

func root(c *cobra.Command, args []string) {

	update, err := c.Flags().GetBool("update")

	if err != nil {
		fmt.Println("Something Went Wrong!", err)
		return
	}

	if !update {
		return
	}

	// TODO: Implement CLI update from Github

}
