package cmd

import (
	"log"

	"github.com/kucicm/tiny-cloud/pkg/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(destroyCmd())
}

func destroyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "delete all resources associated with the profile",
		Long:  "delete all resources associated with the profile",
		Run: func(cmd *cobra.Command, args []string) {
			if err := app.Destroy(); err != nil {
				log.Println(err)
			}
		},
	}
}
