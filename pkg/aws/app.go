package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

type AWS struct {
	ops *tinycloud.Ops
	ecr *ECR
	ecs *ECS
	vm  *VM
	cfg *aws.Config
}

func New() *AWS {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln(err)
	}

	ecr := NewRegistry(cfg)
	ecs := NewECS(cfg)
	vm := NewVM(cfg)

	return &AWS{
		ecr: ecr,
		ecs: ecs,
		vm:  vm,
		cfg: &cfg,
	}
}

func (a *AWS) Init() {

	// vm := newVM(cfg)

	// vm.stop()
}

func (a *AWS) Run(ops tinycloud.Ops) {
	// if destroy?

	// create repo
	a.ecr.createRepository()

	// create bucket

	// create ecs
	a.ecs.create()

	// create vm
	a.vm.Setup()

	// push to repo
	// run on ecs

}

func (a *AWS) Destroy() {
	a.ecr.Destroy()
	a.ecs.destroy()
}
