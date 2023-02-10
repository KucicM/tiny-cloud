package data_test

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
)

func database() (*sql.DB, func()) {
	db := data.SetupDatabes("test.db")
	cleaner := func() {
		data.CloseDatabes()
		os.Remove("test.db")
	}
	return db, cleaner
}

func TestCreateAwsProfile(t *testing.T) {
	db, cleaner := database()
	defer cleaner()

	var count int
	if err := db.QueryRow("SELECT count(*) FROM Profiles").Scan(&count); err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Errorf("expected 0 got %d", count)
	}

	profile := &tinycloud.Profile{
		Name: "test-name-1",
		Settings: &tinycloud.CloudSettings{
			Aws: &tinycloud.AwsSettings{
				AwsRegion:           "region-1",
				AwsAccessKeyId:      "access-key-1",
				AwsSeacretAccessKey: "seacret-acc-key-1",
			},
		},
	}

	if err := data.CreateProfile(profile); err != nil {
		t.Error(err)
	}

	if err := db.QueryRow("SELECT count(*) FROM Profiles").Scan(&count); err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("expected 1 got %d", count)
	}
}

func TestNonUniqueProfileName(t *testing.T) {
	db, cleaner := database()
	defer cleaner()

	var count int
	if err := db.QueryRow("SELECT count(*) FROM Profiles").Scan(&count); err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Errorf("expected 0 got %d", count)
	}

	profile := &tinycloud.Profile{
		Name: "test-name-1",
		Settings: &tinycloud.CloudSettings{
			Aws: &tinycloud.AwsSettings{
				AwsRegion:           "region-1",
				AwsAccessKeyId:      "access-key-1",
				AwsSeacretAccessKey: "seacret-acc-key-1",
			},
		},
	}

	if err := data.CreateProfile(profile); err != nil {
		t.Error(err)
	}

	if err := db.QueryRow("SELECT count(*) FROM Profiles").Scan(&count); err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("expected 1 got %d", count)
	}

	if err := data.CreateProfile(profile); err == nil {
		t.Error("expected error")
	}
}

func TestGetProfilesWithNoProfiles(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	profiles, err := data.GetProfiles()
	if err != nil {
		t.Error(err)
	}

	if len(profiles) != 0 {
		t.Errorf("expected 0 profiles got %d", len(profiles))
	}
}

func TestGetMultipleProfiles(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	profiles := make([]*tinycloud.Profile, 0)
	for i := 0; i < 5; i++ {
		profile := &tinycloud.Profile{
			Name: fmt.Sprintf("test-name-%d", i),
			Settings: &tinycloud.CloudSettings{
				Aws: &tinycloud.AwsSettings{
					AwsRegion:           "region-1",
					AwsAccessKeyId:      "access-key-1",
					AwsSeacretAccessKey: "seacret-acc-key-1",
				},
			},
		}
		profiles = append(profiles, profile)
	}

	expected := make([]*tinycloud.Profile, len(profiles))
	for i, profile := range profiles {
		err := data.CreateProfile(profile)
		if err != nil {
			t.Error(err)
		}

		expected[len(profiles)-i-1] = profile
	}
	expected[0].Active = true

	actual, err := data.GetProfiles()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(tinycloud.Profiles(expected), actual) {
		t.Errorf("expected:\n%+v\nGot:\n%+v", tinycloud.Profiles(expected), actual)
	}
}

func TestGetDefaultActiveProfile(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	var profile *tinycloud.Profile
	for i := 0; i < 5; i++ {
		profile = &tinycloud.Profile{
			Name: fmt.Sprintf("test-name-%d", i),
			Settings: &tinycloud.CloudSettings{
				Aws: &tinycloud.AwsSettings{
					AwsRegion:           "region-1",
					AwsAccessKeyId:      "access-key-1",
					AwsSeacretAccessKey: "seacret-acc-key-1",
				},
			},
		}

		err := data.CreateProfile(profile)
		if err != nil {
			t.Error(err)
		}
	}

	// should return latest profile
	actual, err := data.GetActiveProfile()
	if err != nil {
		t.Error(err)
	}

	profile.Active = true
	if !reflect.DeepEqual(actual, profile) {
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", profile, actual)
	}
}

func TestSetActiveAndGetActive(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	var expected *tinycloud.Profile
	for i := 0; i < 5; i++ {
		profile := &tinycloud.Profile{
			Name: fmt.Sprintf("test-name-%d", i),
			Settings: &tinycloud.CloudSettings{
				Aws: &tinycloud.AwsSettings{
					AwsRegion:           "region-1",
					AwsAccessKeyId:      "access-key-1",
					AwsSeacretAccessKey: "seacret-acc-key-1",
				},
			},
		}

		err := data.CreateProfile(profile)
		if err != nil {
			t.Error(err)
		}

		if i == 2 {
			expected = profile
			expected.Active = true
		}
	}

	if err := data.UpdateProfileToActive(expected.Name); err != nil {
		t.Error(err)
	}

	actual, err := data.GetActiveProfile()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", expected, actual)
	}
}

func TestDeleteNonExistingProfile(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	err := data.DeleteProfile("test")
	if err == nil {
		t.Error("expected error while deleting non existing profile")
	}
}

func TestDeleteExistingProfile(t *testing.T) {
	db, cleaner := database()
	defer cleaner()

	prof := &tinycloud.Profile{
		Name: "test-name-8",
		Settings: &tinycloud.CloudSettings{
			Aws: &tinycloud.AwsSettings{
				AwsRegion:           "region-1",
				AwsAccessKeyId:      "access-key-1",
				AwsSeacretAccessKey: "seacret-acc-key-1",
			},
		},
	}

	err := data.CreateProfile(prof)

	if err != nil {
		t.Error(err)
	}

	var count int
	if err := db.QueryRow("SELECT count(*) FROM v_profiles").Scan(&count); err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("expected 1 got %d", count)
	}

	if err = data.DeleteProfile(prof.Name); err != nil {
		t.Error(err)
	}

	if err = db.QueryRow("SELECT count(*) FROM v_profiles").Scan(&count); err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Errorf("expected 0 got %d", count)
	}
}
