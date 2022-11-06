package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

	profileSetup()

	log.Println("Setting up infra")
	cmd := `terraform -chdir=%s apply -auto-approve -var app=%s -var repo_name=%s`
	out, _ := run(fmt.Sprintf(cmd, path, profile, repoName))
	log.Println(string(out))

	log.Println(img)
	registry.Push(*img)
}

func profileSetup() {
	profileExists := false
	out, _ := run("aws configure list-profiles")
	for _, p := range strings.Split(string(out), "\n") {
		if p == profile {
			profileExists = true
			break
		}
	}

	// setup the profile
	if !profileExists {
		log.Printf("Setup the new user named '%s' on %s\n", profile, "aws TODO")
		cmd := exec.Command("aws", "configure", "--profile", profile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		if err := cmd.Run(); err != nil {
			log.Printf("Failed to run aws config %s\n", err)
			os.Exit(1)
		}
	}

	// set the output
	run(fmt.Sprintf("aws configure set output json --profile %s", profile))
}
