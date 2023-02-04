package state

import (
	"fmt"
	"io"
	"strings"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
	input "github.com/tcnksm/go-input"
)

// menu to create profie, io/out via UI
func CreateNewProfile(in io.Reader, out io.Writer) error {
	ui := &input.UI{
		Writer: out,
		Reader: in,
	}

	var err error

	// name
	name, err := ui.Ask("Name", &input.Options{
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
		Required:  false,
		HideOrder: true,
	})

	if err != nil {
		return err
	}

	// // cloud
	cloud, err := ui.Select("Cloud", tinycloud.SupportedClouds, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		return err
	}

	cloud = strings.ToLower(cloud)
	cloudSettings, err := CreateNewCloudSettings(cloud, ui)
	if err != nil {
		return err
	}

	profile := &tinycloud.Profile{
		Name:        name,
		Description: des,
		Settings:    cloudSettings,
	}

	return data.CreateProfile(profile)
}

// resolves which clouds should be created
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

// aws implementation of settings
func NewAwsCloudSettings(ui *input.UI) (*tinycloud.CloudSettings, error) {
	region, err := ui.Ask("Region", &input.Options{
		Default:   "eu-west-1",
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
	})

	if err != nil {
		return nil, err
	}

	seacretKey, err := ui.Ask("AWS Secret Access Key", &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
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

// list profiles

// set profile to active

// update profile

// delete profile
