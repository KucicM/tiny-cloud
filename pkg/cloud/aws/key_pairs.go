package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type KeyPairApi interface {
	CreateKeyPair(ctx context.Context,
		params *ec2.CreateKeyPairInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateKeyPairOutput, error)
	DeleteKeyPair(ctx context.Context,
		params *ec2.DeleteKeyPairInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteKeyPairOutput, error)
}

func CreateKeyPair(keyName string, tag types.Tag, api KeyPairApi) ([]byte, error) {
	ops := &ec2.CreateKeyPairInput{
		KeyName: aws.String(keyName),
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeKeyPair,
			Tags:         []types.Tag{tag},
		}},
	}
	out, err := api.CreateKeyPair(context.TODO(), ops, opsFn)
	if err != nil {
		return nil, err
	}

	return []byte(*out.KeyMaterial), nil
}

func DeleteKeyPairs(keyNames []string, api KeyPairApi) error {
	for _, keyName := range keyNames {
		ops := &ec2.DeleteKeyPairInput{KeyName: aws.String(keyName)}
		_, err := api.DeleteKeyPair(context.TODO(), ops, opsFn)
		if err != nil {
			return err
		}
	}
	return nil
}
