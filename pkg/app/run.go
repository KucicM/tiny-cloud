package app

import (
	"fmt"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/cloud"
	"github.com/kucicm/tiny-cloud/pkg/state"
)

func Run(req *tinycloud.RunRequest) error {
	profile, err := state.GetActiveProfile()
	if err != nil {
		return err
	}

	// steup vm
	var vm tinycloud.Vm
	switch profile.Settings.ResolveCloudName() {
	case "aws":
		req := cloud.AwsSetupRequest{
			ProfileName:      profile.Name,
			Region:           profile.Settings.AwsRegion,
			AccessKeyId:      profile.Settings.AwsAccessKeyId,
			SeacretAccessKey: profile.Settings.AwsSeacretAccessKey,
		}
		vm, err = cloud.StartAwsVm(req)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("cloud unsupported") // unreachable
	}

	vm.Run(tinycloud.TaskDefinition{})

	return nil
}
