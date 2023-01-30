package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Api interface {
	// DescribeInstances(ctx context.Context,
	// 	params *ec2.DescribeInstancesInput,
	// 	optFns ...func(*ec2.Options)) (ec2.DescribeInstancesOutput, error)
	// RunInstances(ctx context.Context,
	// 	params *ec2.RunInstancesInput,
	// 	optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

type EC2Request struct {
	InstanceType string
}

type EC2Result struct {
	InstanceId string
}

func StartVm(cfg aws.Config, req EC2Request) (EC2Result, error) {
	client := ec2.NewFromConfig(cfg)
	instanceId, ok := checkIfExists(client, req)
	if !ok {
		instanceId, _ = createVm(client, req)
	}
	log.Println("started ", instanceId)
	return EC2Result{}, nil
}

func checkIfExists(client EC2Api, req EC2Request) (string, bool) {
	ops := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{{
			Name:   aws.String("instance-type"),
			Values: []string{req.InstanceType},
		}},
	}

	out, err := client.DescribeInstances(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		log.Fatalf("Error finding instances %s\n", err)
		return "", false
	}

	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			if i.State.Name != types.InstanceStateNameTerminated { // todo add to filter
				return *i.InstanceId, true // todo filter by tag
			}
		}
	}

	return "", false
}

func createVm(client EC2Api, req EC2Request) (string, error) {
	// ops := &ec2.RunInstancesInput{
	// 	MinCount:     aws.Int32(1),
	// 	MaxCount:     aws.Int32(1),
	// 	ImageId:      aws.String("ami-08c41e4d343c2e7ca"), // TODO
	// 	InstanceType: types.InstanceType(req.InstanceType),
	// }

	// if out, err := client.RunInstances(context.TODO(), ops, func(o *ec2.Options) {}); err != nil {
	// 	log.Fatalln(err)
	// 	return "", nil // TODO
	// } else if len(out.Instances) == 1 {
	// 	return *out.Instances[0].InstanceId, nil
	// } else {
	// 	return "", nil // TODO
	// }
	return "", nil
}

func waitStartStatus(instanceId string) bool {
	return false
}

func stopVm(instanceId string) {

}

func DestroyVMs(cfg aws.Config) {
	_ = ec2.NewFromConfig(cfg)
	log.Println("destroy")
}
