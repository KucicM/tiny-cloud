package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
)

// CERATE AWS VM

type AwsSetupRequest struct {
	ProfileName      string
	Region           string
	AccessKeyId      string
	SeacretAccessKey string
	InstanceType     string
	Iam              string
}

func StartVm(req AwsSetupRequest) (*tinycloud.Vm, error) {
	// auth
	client, err := getClient(req.Region, req.AccessKeyId, req.SeacretAccessKey)
	if err != nil {
		return nil, err
	}

	runId, err := data.GetNewRunId(req.ProfileName)
	if err != nil {
		return nil, err
	}

	tag := types.Tag{Key: aws.String("tiny-cloud"), Value: aws.String(runId)}

	if err = CreateSecurityGroup(runId, tag, client); err != nil {
		return nil, err
	}

	if err = AuthorizeSSH(runId, tag, client); err != nil {
		return nil, err
	}

	var sshKey []byte
	if sshKey, err = CreateKeyPair(runId, tag, client); err != nil {
		return nil, err
	}
	data.AddPemKey(runId, sshKey)

	var instanceId string
	if instanceId, err = CretaeEc2(runId, req.InstanceType, req.Iam, tag, client); err != nil {
		return nil, err
	}

	err = waitInstanceStart(instanceId, client)
	if err != nil {
		return nil, err
	}

	var dnsName string
	if dnsName, err = getDNSName(instanceId, client); err != nil {
		return nil, err
	}

	return &tinycloud.Vm{
		Id:      runId,
		SSHKey:  sshKey,
		DNSName: dnsName,
	}, nil
}

// DELETE AWS VM
type AwsDestroyRequest struct {
	ProfileName      string
	Region           string
	AccessKeyId      string
	SeacretAccessKey string
	RunIds           []string
}

// deletes resoures created by the user
func DestroyAws(req AwsDestroyRequest) error {
	// auth
	client, err := getClient(req.Region, req.AccessKeyId, req.SeacretAccessKey)
	if err != nil {
		return err
	}

	if err = DeleteEc2(req.RunIds, client); err != nil {
		return err
	}

	if err = DeleteKeyPairs(req.RunIds, client); err != nil {
		return err
	}

	if err = DeleteSecurityGroups(req.RunIds, client); err != nil {
		return err
	}

	return nil
}

func getClient(region, accessKeyId, seacretAccessKey string) (*ec2.Client, error) {
	creds := credentials.NewStaticCredentialsProvider(
		accessKeyId,
		seacretAccessKey,
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)

	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(cfg)
	return client, err
}

func opsFn(*ec2.Options) {
	// empty
}
