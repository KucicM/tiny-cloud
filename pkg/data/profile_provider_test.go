package data_test

import (
	"database/sql"
	"log"
	"testing"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
)

func database() (*sql.DB, func()) {
	db := data.SetupDatabes("test.db")
	cleaner := func() {
		_, err := db.Exec("DELETE FROM Profiles")
		if err != nil {
			log.Println(err)
		}
	}
	return db, cleaner
}

func TestListNoProfiles(t *testing.T) {
	_, cleaner := database()
	defer cleaner()
	profiles, err := data.ListProfiles()
	if err != nil {
		t.Errorf("did not expect error %s", err)
	}

	if len(profiles) != 0 {
		t.Errorf("did not expect any profiles got %d", len(profiles))
	}
}

func TestAddAndListProfiles(t *testing.T) {
	_, cleaner := database()
	defer cleaner()
	profile := &tinycloud.Profile{Name: "test-profile-1", Description: "test des"}
	err := data.AddProfile(profile)
	if err != nil {
		t.Errorf("did not expect error %s", err)
	}

	profiles, err := data.ListProfiles()
	if err != nil {
		t.Errorf("did not expect error %s", err)
	}

	if len(profiles) != 1 {
		t.Errorf("expected 1 profile got %d", len(profiles))
	}
}
