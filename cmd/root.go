package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolP("debug", "b", false, "print debug states")
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

	// cloud := flag.String("cloud", "aws", "which cloud provider")
	// destroy := flag.Bool("destroy", false, "should delete everything")
	// vmType := flag.String("vm-type", "t2.micro", "vm type to use as ecs")

	// flag.BoolVar(&debug, "debug", false, "debug mode")
	// flag.Parse()

	// var app tinycloud.App
	// switch *cloud {
	// case "aws":
	// 	app = aws.New()
	// default:
	// 	log.Fatalf("no such cloud option %s", *cloud)
	// }

	// if *destroy {
	// 	if err := app.Destroy(); err != nil {
	// 		log.Println(err)
	// 	}
	// 	return
	// }

	// if err := app.Run(tinycloud.Ops{VmType: *vmType}); err != nil {
	// 	log.Fatalln(err)
	// }
}
