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
}

type VmParams struct {
	instanceType string
}

type runningState int

func (s runningState) from(state *types.InstanceState) runningState {
	switch *state.Code {
	case 0:
		return pendingState
	case 16:
		return readyState
	case 32:
		return shuttingDownState
	case 48:
		return terminatedState
	case 80:
		return stoppedState
	}
	return initialState // should never happend
}

const (
	initialState runningState = iota
	createState
	pendingState
	shuttingDownState
	stoppingState
	stoppedState
	readyState
	terminatedState
)

type VmState struct {
	Api          VmAPI
	InstanceType string
	State        runningState
	Id           string
}

func NewVm(api VmAPI, params VmParams) error {
	// vm := &VmState{api, params.instanceType, initial, ""}

	state := initialState

	var instanceId string
	var err error
	for err == nil && state != readyState {
		switch state {
		case initialState:
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
					state = runningState.from(state, instance.State)
				}
			}

			if len(out.Reservations) == 0 {
				state = createState
			}
		case createState:
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
			state = pendingState

		case pendingState, shuttingDownState, stoppingState:
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

			state = runningState.from(state, out.InstanceStatuses[0].InstanceState)
			log.Println("new state ", state)
		case stoppedState:
			// start instance
		case readyState:
			return nil
		case terminatedState:
			state = initialState
		}
	}
	return nil
}
