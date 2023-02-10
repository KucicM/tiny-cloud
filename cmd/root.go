package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
	// rootCmd.PersistentFlags().BoolP("debug", "b", false, "print debug states")
}

var rootCmd = &cobra.Command{
	Use:   "tiny-cloud",
	Short: "Run task on remote VMs",
	Long:  "Run long running task in a cloud or on-prem VMs",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func initConfig() {

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
