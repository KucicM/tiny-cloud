package crud

import (
	"fmt"
	"io"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	input "github.com/tcnksm/go-input"
)

func PrityPrintAllProfiles(writer io.Writer) error {
	profiles, err := GetAllProfiles()
	if err != nil {
		return err
	}

	writer.Write([]byte(profiles.String()))
	return nil
}

func CreateNewProfile(in io.Reader, out io.Writer) error {
	ui := &input.UI{
		Writer: out,
		Reader: in,
	}

	var err error

	// name
	name, err := ui.Ask("Name", &input.Options{
		Default:     "",
		Required:    true,
		Loop:        true,
		HideDefault: false,
		HideOrder:   true,
	})

	if err != nil {
		return err
	}

	// description
	des, err := ui.Ask("Description", &input.Options{
		Default:   "",
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})

	if err != nil {
		return err
	}

	// cloud
	cloud, err := ui.Select("Cloud", tinycloud.SupportedClouds, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return err
	}

	cloudSettings, err := CreateNewCloudSettings(cloud, ui)
	if err != nil {
		return err
	}

	profile := &tinycloud.Profile{
		Name:        name,
		Description: des,
		Cloud:       cloud,
	}

	return Save(profile, cloudSettings)
}

func CreateNewCloudSettings(cloud string, ui *input.UI) (*tinycloud.CloudSettings, error) {
	switch cloud {
	case "aws":
		return NewAwsCloudSettings(ui)
	case "gcp":
	default:
		return nil, fmt.Errorf("cloud %s not supported", cloud)
	}
	return nil, nil
}

func NewAwsCloudSettings(ui *input.UI) (*tinycloud.CloudSettings, error) {
	regions := []string{"eu1"}
	region, err := ui.Select("Region", regions, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return nil, err
	}

	accessKeyId, err := ui.Ask("AWS Access Key ID", &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
		Mask:      true,
	})

	if err != nil {
		return nil, err
	}

	seacretKey, err := ui.Ask("AWS Secret Access Key", &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
		Mask:      true,
	})

	if err != nil {
		return nil, err
	}

	// maybe add default vm type?

	return &tinycloud.CloudSettings{
		AwsRegion:           region,
		AwsAccessKeyId:      accessKeyId,
		AwsSeacretAccessKey: seacretKey,
	}, nil
}
