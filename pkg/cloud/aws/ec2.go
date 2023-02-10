package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

type Ec2 struct {
	Id      string
	SSHKey  []byte
	DNSName string
}

func (e *Ec2) Run(task tinycloud.TaskDefinition) {

}

type Ec2Api interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	DescribeInstanceStatus(ctx context.Context,
		params *ec2.DescribeInstanceStatusInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error)

	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)

	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

func CretaeEc2(runId, instanceType, iam string, api Ec2Api) (string, error) {
	ops := &ec2.RunInstancesInput{
		MinCount:                          aws.Int32(1),
		MaxCount:                          aws.Int32(1),
		ImageId:                           aws.String(iam), // TODO
		InstanceType:                      types.InstanceType(instanceType),
		KeyName:                           aws.String(runId),
		SecurityGroups:                    []string{runId},
		InstanceInitiatedShutdownBehavior: types.ShutdownBehaviorTerminate,
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeInstance,
			Tags: []types.Tag{{
				Key:   aws.String("tiny-cloud"),
				Value: aws.String(runId),
			}},
		}},
	}

	out, err := api.RunInstances(context.TODO(), ops, opsFn)
	if err != nil {
		return "", fmt.Errorf("faild to create new vm with error: %s", err)
	}

	if out == nil {
		return "", fmt.Errorf("run instances return nil")
	}

	if len(out.Instances) == 0 {
		return "", fmt.Errorf("run instances return zero instances")
	}

	instanceId := *out.Instances[0].InstanceId
	if instanceId == "" {
		return "", fmt.Errorf("run instaces return instance without id")
	}

	return instanceId, nil
}

func waitInstanceStart(instanceId string, api Ec2Api) error {
	fmt.Println("waiting instance start")
	ops := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(true),
		InstanceIds:         []string{instanceId},
	}

	// TODO add timeout?
	for {
		out, err := api.DescribeInstanceStatus(context.TODO(), ops, opsFn)

		if err != nil {
			return fmt.Errorf("faild to get vm status %s", err)
		}

		if out == nil || len(out.InstanceStatuses) == 0 {
			continue
		}

		state := out.InstanceStatuses[0].InstanceState
		if state.Name == types.InstanceStateNameRunning {
			break
		}
	}

	return nil
}

func getDNSName(instanceId string, api Ec2Api) (string, error) {
	ops := &ec2.DescribeInstancesInput{InstanceIds: []string{instanceId}}

	out, err := api.DescribeInstances(context.TODO(), ops, opsFn)
	if err != nil {
		return "", fmt.Errorf("cannot get dns name of instanceId: %s, error: %s", instanceId, err)
	}

	if out == nil || len(out.Reservations) == 0 || len(out.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("got unexpected result from describe instances %+v", out)
	}

	instance := out.Reservations[0].Instances[0]
	return *instance.PublicDnsName, nil
}

func DeleteEc2(runIds []string, api Ec2Api) error {
	findOps := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{Name: aws.String("tag-key"), Values: []string{"tiny-cloud"}},
			{Name: aws.String("instance-state-name"), Values: []string{"pending", "running", "stopped"}},
		},
	}

	des, err := api.DescribeInstances(context.TODO(), findOps, opsFn)
	if err != nil {
		return err
	}

	if des == nil || len(des.Reservations) == 0 {
		return nil
	}

	instanceIds := make([]string, 0)
	for _, reservation := range des.Reservations {
		for _, instance := range reservation.Instances {
			if hasRunIdTag(instance.Tags, runIds) {
				instanceIds = append(instanceIds, *instance.InstanceId)
			}
		}
	}

	ops := &ec2.TerminateInstancesInput{InstanceIds: instanceIds}
	if _, err = api.TerminateInstances(context.TODO(), ops, opsFn); err != nil {
		return err
	}

	return waitInstanceTermination(instanceIds, api)
}

func waitInstanceTermination(instanceIds []string, api Ec2Api) error {
	fmt.Println("wating instance termination")
	ops := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(true),
		InstanceIds:         instanceIds,
	}

	// todo timeout?
	anyUp := true
	for anyUp {
		out, err := api.DescribeInstanceStatus(context.TODO(), ops, opsFn)

		if err != nil {
			return fmt.Errorf("faild to get vm status %s", err)
		}

		if out == nil {
			continue
		}

		anyUp = false
		for _, instance := range out.InstanceStatuses {
			if instance.InstanceState.Name != types.InstanceStateNameTerminated {
				anyUp = true
			}
		}
	}
	return nil
}

func hasRunIdTag(tags []types.Tag, runIds []string) bool {
	for _, tag := range tags {
		if *tag.Key == "tiny-cloud" {
			for _, runId := range runIds {
				if runId == *tag.Value {
					return true
				}
			}
		}
	}
	return false
}
