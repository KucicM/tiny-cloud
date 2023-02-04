package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kucicm/tiny-cloud/pkg/state"
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
			// if err := state.PrityPrintAllProfiles(os.Stdout); err != nil {
			// 	fmt.Println(err)
			// 	os.Exit(1)
			// }
		},
	}

	return cmd
}

func profileNewCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "create new profile",
		Long:  "create new profile",
		Run: func(cmd *cobra.Command, args []string) {
			if err := state.CreateNewProfile(os.Stdin, os.Stdout); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	return cmd
}

func profileDeleteCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete",
		Short: "delete profile",
		Long:  "tiny-cloud profile delete <profile-name>",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 || len(strings.TrimSpace(args[0])) == 0 {
				fmt.Println("must provide  profile name")
				return
			}

			// if err := state.DeleteProfile(args[0]); err != nil {
			// 	fmt.Println(err)
			// 	os.Exit(1)
			// }
		},
	}

	return cmd
}
