package cmd

import (
	"fmt"
	"os"

	"github.com/kucicm/tiny-cloud/pkg/crud"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(profileCmd())
}

func profileCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "profile",
		Short: "profile settings",
		Long:  "create/edit profiles",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("profile called")
		},
	}

	cmd.AddCommand(profileListCmd())
	cmd.AddCommand(profileNewCmd())

	return cmd
}

func profileListCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "list existing profiles",
		Long:  "list existing profiles",
		Run: func(cmd *cobra.Command, args []string) {

			profiles, err := crud.ListProfiles()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(profiles)
		},
	}

	return cmd
}

func profileNewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "new new profiles",
		Long:  "create ne profile",
		Run: func(cmd *cobra.Command, args []string) {
			if err := crud.CreateNewProfile(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	return cmd
}
