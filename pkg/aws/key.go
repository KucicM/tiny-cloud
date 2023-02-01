package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type KeyPairApi interface {
	DescribeKeyPairs(ctx context.Context,
		params *ec2.DescribeKeyPairsInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeKeyPairsOutput, error)
}

func NewKeyPair(api KeyPairApi) (struct{}, error) {
	ops := &ec2.DescribeKeyPairsInput{}
	out, err := api.DescribeKeyPairs(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		return struct{}{}, err
	}

	log.Printf("%+v\n", out)
	return struct{}{}, nil
}
