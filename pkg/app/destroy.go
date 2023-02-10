package app

import (
	"fmt"

	"github.com/kucicm/tiny-cloud/pkg/cloud/aws"
	"github.com/kucicm/tiny-cloud/pkg/data"
	"github.com/kucicm/tiny-cloud/pkg/state"
)

func Destroy() error {
	profile, err := state.GetActiveProfile()
	if err != nil {
		return err
	}

	// steup vm
	switch profile.Settings.ResolveCloudName() {
	case "aws":
		runIds, err := data.GetAllRunIds(profile.Name)
		if err != nil {
			return err
		}

		settings := profile.Settings.Aws
		req := aws.AwsDestroyRequest{
			ProfileName:      profile.Name,
			Region:           settings.AwsRegion,
			AccessKeyId:      settings.AwsAccessKeyId,
			SeacretAccessKey: settings.AwsSeacretAccessKey,
			RunIds:           runIds,
		}
		err = aws.DestroyAws(req)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("cloud unsupported") // unreachable
	}

	return nil
}
