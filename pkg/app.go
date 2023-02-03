package tinycloud

import (
	"log"
	"os"
	"path"
)

const Name = "tiny-cloud"
const DbName = "tiny-cloud.db"

var ConfigPath string
var DbPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("$HOME not defined")
	}

	ConfigPath = path.Join(home, ".tiny-cloud")
	log.Println(ConfigPath)
	err = os.MkdirAll(ConfigPath, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}

	DbPath = path.Join(ConfigPath, DbName)
}

type Ops struct {
	VmType string
}

type App interface {
	Run(ops Ops) error
	Destroy() error
}
