package main

import (
	"flag"
	"fmt"
	"log"
)

const path = "terraform"
const profile = "tiny-cloud"
const repoName = profile + "-repo"

var debug = false

func main() {

	destroy := flag.Bool("destroy", false, "should delete everything")
	img := flag.String("image", "", "image name to run")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	run("terraform init")
	registry := &ECR{}

	if *destroy {
		log.Println("Destroy!!")
		registry.DeleteImgs()
		out, _ := run(fmt.Sprintf("terraform -chdir=%s destroy -auto-approve", path))
		log.Println(string(out))
		return
	}

	log.Println("Setting up infra")
	cmd := `terraform -chdir=%s apply -auto-approve -var app=%s -var repo_name=%s`
	out, _ := run(fmt.Sprintf(cmd, path, profile, repoName))
	log.Println(string(out))

	log.Println(img)
	registry.Push(*img)
}
