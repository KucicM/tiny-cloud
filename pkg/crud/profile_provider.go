package crud

import (
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

	profile := &tinycloud.Profile{}
	var err error

	// name
	profile.Name, err = ui.Ask("Name", &input.Options{
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
	profile.Description, err = ui.Ask("Description", &input.Options{
		Default:   "",
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})

	if err != nil {
		return err
	}

	// cloud
	return SaveProfile(profile)
}
