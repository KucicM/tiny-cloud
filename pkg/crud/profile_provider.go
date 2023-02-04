package crud

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

func ListProfiles() (tinycloud.Profiles, error) {
	return GetAllProfiles()
}

func CreateNewProfile() error {
	fmt.Println("Create new profile")
	profile := &tinycloud.Profile{}
	if err := setStrInput("Name: ", &profile.Name); err != nil {
		return err
	}

	if err := setStrInput("Description: ", &profile.Description); err != nil {
		return err
	}

	// name
	// description
	// cloud
	return SaveProfile(profile)
}

func setStrInput(prompt string, f *string) error {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err == nil {
		*f = strings.TrimSuffix(in, "\n")
	}
	return err

}
