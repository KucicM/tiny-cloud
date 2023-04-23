package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

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

func CreateEC2(runId, instanceType, iam string, tag types.Tag, api Ec2Api) (string, error) {
	ops := &ec2.RunInstancesInput{
		MinCount:                          aws.Int32(1),
		MaxCount:                          aws.Int32(1),
		ImageId:                           aws.String(iam),
		InstanceType:                      types.InstanceType(instanceType),
		KeyName:                           aws.String(runId),
		SecurityGroups:                    []string{runId},
		InstanceInitiatedShutdownBehavior: types.ShutdownBehaviorTerminate,
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeInstance,
			Tags:         []types.Tag{tag},
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
	doneCondition := func(in *ec2.DescribeInstanceStatusOutput) bool {
		state := in.InstanceStatuses[0].InstanceState
		return state.Name == types.InstanceStateNameRunning
	}
	return waitInstanceStatus([]string{instanceId}, doneCondition, api)
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
			{Name: aws.String("tag:tiny-cloud"), Values: runIds},
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
			instanceIds = append(instanceIds, *instance.InstanceId)
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

	doneCondition := func(in *ec2.DescribeInstanceStatusOutput) bool {
		for _, instance := range in.InstanceStatuses {
			if instance.InstanceState.Name != types.InstanceStateNameTerminated {
				return false
			}
		}
		return true
	}

	return waitInstanceStatus(instanceIds, doneCondition, api)
}

func waitInstanceStatus(instanceIds []string,
	doneCond func(*ec2.DescribeInstanceStatusOutput) bool, api Ec2Api) error {

	ops := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(true),
		InstanceIds:         instanceIds,
	}

	for {
		out, err := api.DescribeInstanceStatus(context.TODO(), ops, opsFn)

		if err != nil {
			return fmt.Errorf("faild to get vm status %s", err)
		}

		if out == nil || len(out.InstanceStatuses) == 0 {
			continue
		}

		if doneCond(out) {
			break
		}
	}
	return nil
}
