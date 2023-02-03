package tinycloud_test

import (
	"testing"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
)

// TODO clear dbs and create separete db for tests

func TestListNoProfiles(t *testing.T) {
	profiles, err := tinycloud.ListProfiles()
	if err != nil {
		t.Errorf("did not expect error %s", err)
	}

	if len(profiles) != 0 {
		t.Errorf("did not expect any profiles got %d", len(profiles))
	}
}

func TestAddAndListProfiles(t *testing.T) {
	profile := &tinycloud.Profile{Name: "test-profile-1", Description: "test des"}
	err := tinycloud.AddProfile(profile)
	if err != nil {
		t.Errorf("did not expect error %s", err)
	}

	profiles, err := tinycloud.ListProfiles()
	if err != nil {
		t.Errorf("did not expect error %s", err)
	}

	if len(profiles) != 1 {
		t.Errorf("expected 1 profile got %d", len(profiles))
	}
}
