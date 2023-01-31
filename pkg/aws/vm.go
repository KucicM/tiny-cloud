package aws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go/aws"
)

type VmAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)
	DescribeInstanceStatus(ctx context.Context,
		params *ec2.DescribeInstanceStatusInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceStatusOutput, error)
	StartInstances(ctx context.Context,
		params *ec2.StartInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
}

type VmRequest struct {
	vmType string
	debug  bool
}

func NewVm(api VmAPI, req VmRequest) error {
	state := Initial

	var vmId string
	var err error
	for err == nil && state != Ready {
		prevState := state
		switch state {
		case Initial:
			vmId, state, err = findVmIfExists(api, req.vmType)
		case Create:
			vmId, state, err = createVm(api, req.vmType)
		case Pending, ShuttingDown, Stopping:
			time.Sleep(time.Second * 1) // wait for next state
			state, err = getVmState(api, vmId)
		case Stopped:
			state, err = startVm(api, vmId)
		case Ready:
			return nil
		case Terminated:
			state = Initial
		}
		if req.debug {
			debug(prevState, state, vmId, req.vmType, err)
		}
	}
	return nil
}

func findVmIfExists(api VmAPI, vmType string) (string, VmState, error) {
	// todo add filter by tag
	ops := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{{
			Name:   aws.String("instance-type"),
			Values: []string{vmType},
		}, {
			Name: aws.String("instance-state-name"),
			Values: []string{"pending", "running",
				"shutting-down", "stopping", "stopped"},
		}},
	}

	out, err := api.DescribeInstances(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		return "", Error, fmt.Errorf("failed on describe vm with error: %s", err)
	}

	for _, r := range out.Reservations {
		for _, instance := range r.Instances {
			return *instance.InstanceId, VmState.fromAwsState(Initial, instance.State), nil
		}
	}

	return "", Create, nil
}

func createVm(api VmAPI, vmType string) (string, VmState, error) {
	ops := &ec2.RunInstancesInput{
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		ImageId:      aws.String("ami-08c41e4d343c2e7ca"), // TODO
		InstanceType: types.InstanceType(vmType),
	}
	out, err := api.RunInstances(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		return "", Error, fmt.Errorf("faild to create new vm with error: %s", err)
	}

	if len(out.Instances) == 0 {
		return "", Error, fmt.Errorf("got zero vm after create")
	}
	return *out.Instances[0].InstanceId, Pending, nil
}

func getVmState(api VmAPI, vmId string) (VmState, error) {
	if vmId == "" {
		return Error, fmt.Errorf("cannot check state without vm id")
	}
	ops := &ec2.DescribeInstanceStatusInput{
		IncludeAllInstances: aws.Bool(true),
		InstanceIds:         []string{vmId},
	}

	out, err := api.DescribeInstanceStatus(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		return Error, fmt.Errorf("faild to get vm status %s", err)
	}

	if len(out.InstanceStatuses) == 0 {
		return Error, fmt.Errorf("got zero vm state with id %s", vmId)
	}
	return VmState.fromAwsState(Pending, out.InstanceStatuses[0].InstanceState), nil
}

func startVm(api VmAPI, vmId string) (VmState, error) {
	if vmId == "" {
		return Error, fmt.Errorf("cannot start vm without id")
	}

	ops := &ec2.StartInstancesInput{InstanceIds: []string{vmId}}
	_, err := api.StartInstances(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		return Error, fmt.Errorf("faild to start instance %s", err)
	}
	return Pending, nil
}

func debug(pState, state VmState, vmId, vmType string, err error) {
	var errStr string
	if err != nil {
		errStr = fmt.Sprintf(" error: %s", err)
	}

	if vmId != "" {
		vmId = fmt.Sprintf(" vmId: %s", vmId)
	}
	log.Printf("%s -> %s: type: %s%s%s\n", pState, state, vmType, vmId, errStr)
}
