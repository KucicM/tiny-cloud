package main

import (
	"flag"
	"log"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/aws"
)

var debug = false

func main() {

	cloud := flag.String("cloud", "aws", "which cloud provider")
	destroy := flag.Bool("destroy", false, "should delete everything")
	vmType := flag.String("vm-type", "t2.micro", "vm type to use as ecs")

	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	var app tinycloud.App
	switch *cloud {
	case "aws":
		app = aws.New()
	default:
		log.Fatalf("no such cloud option %s", *cloud)
	}

	if *destroy {
		if err := app.Destroy(); err != nil {
			log.Println(err)
		}
		return
	}

	if err := app.Run(tinycloud.Ops{VmType: *vmType}); err != nil {
		log.Fatalln(err)
	}
}
