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

	img := flag.String("image", "test", "image name to run")
	vm := flag.String("vm-type", "t2.micro", "vm type to use as ecs")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()


	// make it simpler
	// use std not loggers (loggers maybe to file?)

	var app tinycloud.App
	// todo convert to lowercase
	switch *cloud {
	case "aws":
		app = aws.New()
	default:
		log.Fatalf("%s not supported", *cloud)
		return
	}

	app.Init()

	if *destroy {
		log.Println("Destroy!!")
		app.Destroy()
		return
	}

	ops := tinycloud.Ops{
		Image:  *img,
		VmType: *vm,
	}

	app.Run(ops)
}
