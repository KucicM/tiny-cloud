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
	ec2Client := ec2.NewFromConfig(a.cfg)

	op := &ec2.DescribeKeyPairsInput{KeyNames: []string{tinycloud.PROFILE_NAME}}
	out, err := ec2Client.DescribeKeyPairs(context.TODO(), op, func(o *ec2.Options) {})
	if err != nil {
		return err
	}

	if len(out.KeyPairs) == 0 {
		// todo create key
		return nil
	}

	info := out.KeyPairs[0]
	name := info.KeyName

	log.Printf("out %v\n", len(out.KeyPairs))
	log.Printf("out %v\n", *name)

	// ec2Client.DescribeKeyPairs(ec2Client)

	// key, err := NewKeyPair(ec2Client)

	NewVm(ec2Client, VmRequest{ops.VmType, true, tinycloud.KEY_NAME})

	// todo push docker image

	// clean docker staging image

	return nil
}

func (a *AWS) Destroy() error {
	// DestroyVMs(a.cfg)
	// keys
	return nil
}
