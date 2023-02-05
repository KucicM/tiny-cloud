package cmd

import (
	"fmt"
	"os"

	"github.com/kucicm/tiny-cloud/pkg/state"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newProfileCmd())
	rootCmd.AddCommand(useProfileCmd())
	rootCmd.AddCommand(listProfilesCmd())
	rootCmd.AddCommand(deleteProfileCmd())
}

func newProfileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "new",
		Short: "Create new profile",
		Long:  "Create new profile",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				fmt.Println("Error: expect no argumets")
				os.Exit(1)
			}
			if err := state.CreateNewProfile(os.Stdin, os.Stdout); err != nil {
				fmt.Printf("Error: %s", err)
			}
		},
	}
}

func useProfileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use",
		Short: "Select profile to use",
		Long:  "Select profile to use",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("Error: requires exactly 1 profile, recived %d", len(args))
				os.Exit(1)
			}

			profileName := args[0]
			if err := state.SetActive(profileName); err != nil {
				fmt.Printf("Error: %s", err)
				os.Exit(1)
			}
		},
	}
}

func listProfilesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available profiles",
		Long:  "List available profiles",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				fmt.Println("expected no args")
				os.Exit(1)
			}

			profiles, err := state.ListProfiles()
			if err != nil {
				fmt.Printf("Error: %s", err)
				os.Exit(1)
			}
			fmt.Print(profiles.String())
		},
	}
}

func deleteProfileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete profile(s)",
		Long:  "Delete profile(s)",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Printf("Eror: expecetd at least 1 argument")
				os.Exit(1)
			}

			state.DeleteProfile(args...)
		},
	}
}
