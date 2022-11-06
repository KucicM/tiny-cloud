package main

import (
	"log"
	"os/exec"
)

func run(run string) ([]byte, error) {
	if debug {
		log.Println("Running", run)
	}
	return exec.Command("bash", "-c", run).CombinedOutput()
}
