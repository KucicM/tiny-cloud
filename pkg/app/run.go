package app

import (
	"fmt"
	"log"

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
		settings := profile.Settings.Aws
		req := aws.AwsSetupRequest{
			ProfileName:      profile.Name,
			Region:           settings.AwsRegion,
			AccessKeyId:      settings.AwsAccessKeyId,
			SeacretAccessKey: settings.AwsSeacretAccessKey,
			InstanceType:     req.VmType,
			Iam:              "ami-06c39ed6b42908a36", // todo from db defaults
		}

        bucketName, err := aws.CreateS3(req)
        if err != nil {
			return err
		}
        req.BucketName = bucketName
        log.Println(bucketName)

		vm, err = aws.StartVm(req)
		if err != nil {
			return err
		}

        log.Println("DONE")


	default:
		return fmt.Errorf("cloud unsupported") // unreachable
	}

	return cloud.Run(tinycloud.TaskDefinition{
		SSHKey:        vm.SSHKey,
		DNSName:       vm.DNSName,
		DockerImageId: req.DockerImageId,
	})
}

