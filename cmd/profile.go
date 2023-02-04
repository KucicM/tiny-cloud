package cmd

import (
	"fmt"
	"os"

	"github.com/kucicm/tiny-cloud/pkg/data"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "profile settings",
	Long:  "create/edit profiles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profile called")
	},
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "list existing profiles",
	Long:  "list existing profiles",
	Run: func(cmd *cobra.Command, args []string) {

		profiles, err := data.ListProfiles()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(profiles)
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
