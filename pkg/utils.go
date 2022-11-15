package tinycloud

import (
	"log"
	"os/exec"
)

func Run(run string) ([]byte, error) {
	log.Println("Running", run)
	return exec.Command("bash", "-c", run).CombinedOutput()
}
