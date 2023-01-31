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

// state enmu

type VmState int64

const (
	Initial VmState = iota
	Create
	Pending
	ShuttingDown
	Stopping
	Stopped
	Ready
	Terminated
	Error
)

func (s VmState) String() string {
	strs := []string{"INITIAL", "CREATE", "PENDING", "SHUTTINGDOWN",
		"STOPPING", "STOPPED", "READY", "TERMINATED", "ERROR"}
	if s < Initial || s > Error {
		return "UNKNOWN"
	}
	return strs[s]
}

func (s VmState) fromAwsState(state *types.InstanceState) VmState {
	switch *state.Code {
	case 0:
		return Pending
	case 16:
		return Ready
	case 32:
		return ShuttingDown
	case 48:
		return Terminated
	case 64:
		return Stopping
	case 80:
		return Stopped
	default:
		log.Printf("got unexpected aws state %+v\n", state)
		return Initial
	}
}

// done state enum

type VmParams struct {
	instanceType string
}

func NewVm(api VmAPI, params VmParams) error {
	state := Initial

	var instanceId string
	var err error
	for err == nil && state != Ready {
		switch state {
		case Initial:
			log.Println("check if exists")
			// todo add filter by tag
			ops := &ec2.DescribeInstancesInput{
				Filters: []types.Filter{{
					Name:   aws.String("instance-type"),
					Values: []string{params.instanceType},
				}, {
					Name: aws.String("instance-state-name"),
					Values: []string{"pending", "running",
						"shutting-down", "stopping", "stopped"},
				}},
			}

			out, err := api.DescribeInstances(context.TODO(), ops, func(o *ec2.Options) {})
			if err != nil {
				return err
			}

			for _, r := range out.Reservations {
				for _, instance := range r.Instances {
					instanceId = *instance.InstanceId
					state = VmState.fromAwsState(state, instance.State)
				}
			}

			if len(out.Reservations) == 0 {
				state = Create
			}
		case Create:
			log.Println("create vm")

			ops := &ec2.RunInstancesInput{
				MinCount:     aws.Int32(1),
				MaxCount:     aws.Int32(1),
				ImageId:      aws.String("ami-08c41e4d343c2e7ca"), // TODO
				InstanceType: types.InstanceType(params.instanceType),
			}
			out, err := api.RunInstances(context.TODO(), ops, func(o *ec2.Options) {})
			if err != nil {
				return err
			}

			if len(out.Instances) != 1 {
				return fmt.Errorf("unexpected number of instances, %d", len(out.Instances))
			}
			instanceId = *out.Instances[0].InstanceId
			state = Pending
		case Pending, ShuttingDown, Stopping:
			time.Sleep(time.Second * 5) // wait for next state

			ops := &ec2.DescribeInstanceStatusInput{
				IncludeAllInstances: aws.Bool(true),
				InstanceIds:         []string{instanceId},
			}

			out, err := api.DescribeInstanceStatus(context.TODO(), ops, func(o *ec2.Options) {})
			if err != nil {
				return err
			}

			if len(out.InstanceStatuses) != 1 {
				return fmt.Errorf("unexpected count of instance status %d", len(out.InstanceStatuses))
			}

			state = VmState.fromAwsState(state, out.InstanceStatuses[0].InstanceState)
			log.Println("new state ", state)
		case Stopped:
			log.Println("start vm")
			ops := &ec2.StartInstancesInput{
				InstanceIds: []string{instanceId},
			}

			_, err := api.StartInstances(context.TODO(), ops, func(o *ec2.Options) {})
			if err != nil {
				return err
			}

			state = Pending
		case Ready:
			return nil
		case Terminated:
			state = Initial
		}
	}
	return nil
}
