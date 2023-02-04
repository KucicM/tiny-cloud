package crud_test

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/crud"
)

func database() (*sql.DB, func()) {
	db := crud.SetupDatabes("test.db")
	cleaner := func() {
		crud.CloseDatabes()
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
			AwsRegion:           "region-1",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	if err := crud.CreateProfile(profile); err != nil {
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
			AwsRegion:           "region-1",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	if err := crud.CreateProfile(profile); err != nil {
		t.Error(err)
	}

	if err := db.QueryRow("SELECT count(*) FROM Profiles").Scan(&count); err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Errorf("expected 1 got %d", count)
	}

	if err := crud.CreateProfile(profile); err == nil {
		t.Error("expected error")
	}
}

func TestGetProfilesWithNoProfiles(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	profiles, err := crud.GetProfiles()
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
				AwsRegion:           "region-1",
				AwsAccessKeyId:      "access-key-1",
				AwsSeacretAccessKey: "seacret-acc-key-1",
			},
		}
		profiles = append(profiles, profile)
	}

	for _, profile := range profiles {
		err := crud.CreateProfile(profile)
		if err != nil {
			t.Error(err)
		}
	}

	actual, err := crud.GetProfiles()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(tinycloud.Profiles(profiles), actual) {
		t.Errorf("expected:\n%+v\nGot:\n%+v", tinycloud.Profiles(profiles), actual)
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
				AwsRegion:           "region-1",
				AwsAccessKeyId:      "access-key-1",
				AwsSeacretAccessKey: "seacret-acc-key-1",
			},
		}

		err := crud.CreateProfile(profile)
		if err != nil {
			t.Error(err)
		}
	}

	// should return latest profile
	actual, err := crud.GetActiveProfile()
	if err != nil {
		t.Error(err)
	}

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
				AwsRegion:           "region-1",
				AwsAccessKeyId:      "access-key-1",
				AwsSeacretAccessKey: "seacret-acc-key-1",
			},
		}

		err := crud.CreateProfile(profile)
		if err != nil {
			t.Error(err)
		}

		if i == 2 {
			expected = profile
		}
	}

	if err := crud.UpdateProfileToActive(expected.Name); err != nil {
		t.Error(err)
	}

	actual, err := crud.GetActiveProfile()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", expected, actual)
	}
}

func TestUpdateNonExistingProfile(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	oldP := &tinycloud.Profile{
		Name: "test-name-8",
		Settings: &tinycloud.CloudSettings{
			AwsRegion:           "region-1",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	newP := &tinycloud.Profile{
		Name: "test-name-8",
		Settings: &tinycloud.CloudSettings{
			AwsRegion:           "region-1",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	if err := crud.UpdateProfile(oldP, newP); err == nil {
		t.Error("expected error")
	}

}

func TestUpdateExistingProfile(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	oldP := &tinycloud.Profile{
		Name: "test-name-8",
		Settings: &tinycloud.CloudSettings{
			AwsRegion:           "region-1",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	newP := &tinycloud.Profile{
		Name: "test-name-8",
		Settings: &tinycloud.CloudSettings{
			AwsRegion:           "region0",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	// create
	if err := crud.CreateProfile(oldP); err != nil {
		t.Error(err)
	}

	// update
	if err := crud.UpdateProfile(oldP, newP); err != nil {
		t.Error(err)
	}

	// get
	acutal, err := crud.GetActiveProfile()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(newP, acutal) {
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", newP, acutal)
		t.Errorf("expected:\n%+v\ngot:\n%+v\n", newP.Settings, acutal.Settings)
	}

}

func TestDeleteNonExistingProfile(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	err := crud.DeleteProfile("test")
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
			AwsRegion:           "region-1",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	err := crud.CreateProfile(prof)

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

	if err = crud.DeleteProfile(prof.Name); err != nil {
		t.Error(err)
	}

	if err = db.QueryRow("SELECT count(*) FROM v_profiles").Scan(&count); err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Errorf("expected 0 got %d", count)
	}
}
