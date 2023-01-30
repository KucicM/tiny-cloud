package aws_test

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	// tinycloud "github.com/kucicm/tiny-cloud/pkg/aws"
)

type describeInstancesResponse struct {
	value ec2.DescribeInstancesOutput
	err   error
}

type describeInstancesState struct {
	response           describeInstancesResponse
	expectedParams     *ec2.DescribeInstancesInput
	expectedToBeCalled bool
	expectedOrderCall  int
}

type mockEC2Api struct {
	callCount int
	diState   *describeInstancesState
}

func (m *mockEC2Api) DescribeInstances(ctx context.Context,
	params *ec2.DescribeInstancesInput,
	optFns ...func(*ec2.Options)) (ec2.DescribeInstancesOutput, error) {
	m.callCount++

	state := m.diState
	if state == nil {
		log.Fatalln("DescribeInstances state not state")
	}
	if !state.expectedToBeCalled {
		log.Fatalln("DescribeInstances but it should not")
	}

	if m.callCount != state.expectedOrderCall {
		log.Fatalf(
			"DescribeInstances call order should be %d but is %d\n",
			state.expectedOrderCall, m.callCount)
	}

	if !reflect.DeepEqual(params, state.expectedParams) {
		log.Fatalf("DescribeInstances expected params %+v got %+v\n", params, state.expectedParams)
	}

	return m.diState.response.value, m.diState.response.err
}

func (m *mockEC2Api) RunInstances(ctx context.Context,
	params *ec2.RunInstancesInput,
	optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error) {

	return nil, nil
}

func TestVmDoesNotExistCreateNew(t *testing.T) {
	// mock := &mockEC2Api{
	// 	diState: &describeInstancesState{
	// 		expectedToBeCalled: true,
	// 		expectedOrderCall:  1,
	// 		expectedParams: &ec2.DescribeInstancesInput{
	// 			Filters: []types.Filter{{
	// 				Name:   aws.String("instance-type"),
	// 				Values: []string{"t2.micro"},
	// 			}},
	// 		},
	// 	},
	// }
	// tinycloud.AddVmToCluster(mock, types.InstanceTypeT2Micro)
}
