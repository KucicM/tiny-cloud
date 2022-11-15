package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

type AWS struct {
	ops *tinycloud.Ops
	ecr *ECR
}

func New() *AWS {
	return &AWS{}
}

func (a *AWS) Init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	vm := newVM(cfg)

	vm.stop()
}

func (a *AWS) Run(ops tinycloud.Ops) {
	// a.ecr.createRepository()

}

func (a *AWS) Destroy() {
}
