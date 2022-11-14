package main

import (
	"context"
	"flag"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
)

var debug = false

func main() {

	destroy := flag.Bool("destroy", false, "should delete everything")
	img := flag.String("image", "", "image name to run")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	repo := NewRegistry(cfg)

	if *destroy {
		log.Println("Destroy!!")
		repo.Destroy()
		return
	}

	repo.Push(*img)
}
