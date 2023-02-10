package cmd

import (
	"log"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd())
}

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a task in a cloud",
		Long:  "Run a task in a cloud",
		Run: func(cmd *cobra.Command, args []string) {
			req := &tinycloud.RunRequest{}
			req.LocalFilesPath, _ = cmd.Flags().GetString("src-path")
			req.VmType, _ = cmd.Flags().GetString("vm-type")
			req.DataOutPath, _ = cmd.Flags().GetString("data-out-path")
			if err := app.Run(req); err != nil {
				log.Println(err)
			}
		},
	}

	cmd.Flags().String(
		"src-path",
		"./*",
		"which file(s) should be copied to vm",
	)
	cmd.MarkFlagRequired("src-path")

	cmd.Flags().String(
		"vm-type",
		"",
		"which type of vm should be used",
	)
	cmd.MarkFlagRequired("vm-type")

	cmd.Flags().String(
		"data-out-path",
		"",
		"path which should be saved after run",
	)

	return cmd
}
