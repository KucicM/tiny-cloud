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

		req := aws.AwsDestroyRequest{
			ProfileName:      profile.Name,
			Region:           profile.Settings.AwsRegion,
			AccessKeyId:      profile.Settings.AwsAccessKeyId,
			SeacretAccessKey: profile.Settings.AwsSeacretAccessKey,
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
