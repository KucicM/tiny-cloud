package aws

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

type AWS struct {
	ops *tinycloud.Ops
	cfg aws.Config
}

func New() *AWS {
	//cfg, err := config.LoadDefaultConfig(context.TODO())
	creds := credentials.NewEnvCredentials()
	v, err := creds.Get()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v\n", v)

	return &AWS{}
}

func (a *AWS) Run(ops tinycloud.Ops) error {
	// prepare docker image

	vmReq := EC2Request{InstanceType: ops.VmType}
	_, err := StartVm(a.cfg, vmReq)
	if err != nil {
		return err
	}

	// todo push docker image

	// clean docker staging image

	return nil
}

func (a *AWS) Destroy() error {
	DestroyVMs(a.cfg)
	// keys
	return nil
}
