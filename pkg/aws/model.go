package aws

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

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
