package aws

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type VM struct {
	client     *ec2.Client
	instanceId string
}

// TODO add vm options
func NewVM(cfg aws.Config) *VM {
	return &VM{client: ec2.NewFromConfig(cfg)}
	// e.setup()
	// return e
}

func (e *VM) Setup() {
	e.init()
	e.awit()
	log.Printf("setup done instanceId: %s\n", e.instanceId)
}

func (e *VM) awit() {
	// wait to start
	status, err := e.getStatus()
	for err == nil && status != types.InstanceStateNameRunning {
		if status == types.InstanceStateNameStopped {
			e.start()
		}
		log.Println("waiting to start machine")
		time.Sleep(time.Second * 15)
		status, err = e.getStatus()
	}
}

func (e *VM) start() {
	log.Println("Starting vm")
	ops := &ec2.StartInstancesInput{InstanceIds: []string{e.instanceId}}
	_, err := e.client.StartInstances(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		log.Fatalf("Failed to start vm %s\n", err)
	}
}

func (e *VM) stop() {
	log.Println("Stopping vm")
	ops := &ec2.StopInstancesInput{
		InstanceIds: []string{e.instanceId},
		Force:       aws.Bool(true),
	}
	_, err := e.client.StopInstances(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		log.Fatalf("Failed to stop VM %s\n", err)
	}
}

func (e *VM) init() {
	// todo check if exists and create
	t := types.InstanceTypeT2Micro
	e.resolveInstanceId(t)
}

func (e *VM) resolveInstanceId(t types.InstanceType) {
	instanceId := e.findExistingVm(t)
	if instanceId == "" {
		instanceId = e.create()
	}

	e.instanceId = instanceId
	if instanceId == "" {
		log.Fatalln("Faild to create VM")
	}
}

func (e *VM) findExistingVm(t types.InstanceType) string {
	ops := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{{
			Name:   aws.String("instance-type"),
			Values: []string{string(t)},
		}},
	}
	out, err := e.client.DescribeInstances(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		log.Fatalf("Error finding instances %s\n", err)
	}
	for _, r := range out.Reservations {
		for _, i := range r.Instances {
			if i.State.Name != types.InstanceStateNameTerminated {
				return *i.InstanceId // todo filter by tag
			}
		}
	}
	return ""
}

func (e *VM) create() string {
	// create iam role
	iam.CreateRoleInput{}

	// TODO install rex and params
	userData := `#!/bin/bash echo ECS_CLUSTER=tiny-cloud-ecs >> /etc/ecs/ecs.config`
	bs64 := base64.StdEncoding.EncodeToString([]byte(userData))

	ops := &ec2.RunInstancesInput{
		MaxCount: aws.Int32(1),
		MinCount: aws.Int32(1),
		ImageId:  aws.String("ami-08c41e4d343c2e7ca"), // ecs optimized
		// IamInstanceProfile: //TODO create IAM role just for this or use tiny-cloud?
		InstanceType: types.InstanceTypeT2Micro,
		UserData:     aws.String(bs64),
	}
	out, err := e.client.RunInstances(context.TODO(), ops, func(o *ec2.Options) {})
	if err != nil {
		log.Println(err)
		return ""
	}
	return *out.Instances[0].InstanceId
}

func (e *VM) getStatus() (types.InstanceStateName, error) {
	out, err := e.client.DescribeInstanceStatus(
		context.TODO(),
		&ec2.DescribeInstanceStatusInput{
			IncludeAllInstances: aws.Bool(true),
			InstanceIds:         []string{e.instanceId},
		},
		func(o *ec2.Options) {},
	)
	if err != nil {
		return types.InstanceStateNamePending, err
	}
	return out.InstanceStatuses[0].InstanceState.Name, nil
}
