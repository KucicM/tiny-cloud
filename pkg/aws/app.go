package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

type AWS struct {
	ops *tinycloud.Ops
	cfg aws.Config
}

func New() *AWS {
	// maybe add cofigurable profile?

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(tinycloud.PROFILE_NAME),
	)

	if err != nil {
		log.Fatalln(err)
	}

	return &AWS{cfg: cfg}
}

func (a *AWS) Run(ops tinycloud.Ops) error {
	// prepare docker image

	// vmReq := EC2Request{InstanceType: ops.VmType}
	client := ec2.NewFromConfig(a.cfg)
	NewVm(client, VmRequest{ops.VmType, true})

	// todo push docker image

	// clean docker staging image

	return nil
}

func (a *AWS) Destroy() error {
	// DestroyVMs(a.cfg)
	// keys
	return nil
}
