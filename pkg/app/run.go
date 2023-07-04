package app

import (
	"fmt"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/cloud"
	"github.com/kucicm/tiny-cloud/pkg/cloud/aws"
	"github.com/kucicm/tiny-cloud/pkg/state"
)

func Run(req *tinycloud.RunRequest) error {
	profile, err := state.GetActiveProfile()
	if err != nil {
		return err
	}

	found, err := cloud.DoseDockerImageExists(req.DockerImageId)
	if err != nil {
		return err
	}

	if !found {
		return fmt.Errorf("cannot find image with Id='%s'", req.DockerImageId)
	}

	// steup vm
	var vm *tinycloud.Vm
	switch profile.Settings.ResolveCloudName() {
	case "aws":

        // create VMw
		settings := profile.Settings.Aws
		req := aws.AwsSetupRequest{
			ProfileName:      profile.Name,
			Region:           settings.AwsRegion,
			AccessKeyId:      settings.AwsAccessKeyId,
			SeacretAccessKey: settings.AwsSeacretAccessKey,
			InstanceType:     req.VmType,
			Iam:              "ami-06c39ed6b42908a36", // todo from db defaults
		}
        /*
		vm, err = aws.StartVm(req)
		if err != nil {
			return err
		}
        */

        if _, err := aws.CreateS3(req); err != nil {
			return err
		}
	default:
		return fmt.Errorf("cloud unsupported") // unreachable
	}

	return cloud.Run(tinycloud.TaskDefinition{
		SSHKey:        vm.SSHKey,
		DNSName:       vm.DNSName,
		DockerImageId: req.DockerImageId,
	})
}

