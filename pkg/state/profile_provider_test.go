package state_test

import (
	"bytes"
	"database/sql"
	"os"
	"strings"
	"testing"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
	"github.com/kucicm/tiny-cloud/pkg/state"
)

func database() (*sql.DB, func()) {
	db := data.SetupDatabes("test.db")
	cleaner := func() {
		data.CloseDatabes()
		os.Remove("test.db")
	}
	return db, cleaner
}

func TestCreateNewProfile(t *testing.T) {
	db, cleaner := database()
	defer cleaner()

	var count int
	if err := db.QueryRow("SELECT count(*) FROM v_profiles;").Scan(&count); err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Errorf("expected count 0 got %d", count)
	}

	in := &bytes.Buffer{}
	in.Write([]byte("test-name\n"))
	in.Write([]byte("\n"))
	in.Write([]byte("1\n")) // aws
	in.Write([]byte("\n"))  // use default
	in.Write([]byte("xxxx-xxxx-xxxx\n"))
	in.Write([]byte("ffff-ffff-ffff\n"))
	out := &bytes.Buffer{}
	if err := state.CreateNewProfile(in, out); err != nil {
		t.Error(err)
	}

	expected := []string{
		"Name:",
		"Description:",
		"Cloud",
		"",
		"1. aws",
		"2. gcp",
		"",
		"Enter a number:",
		"Region (Default is eu-west-1):",
		"AWS Access Key ID:",
		"AWS Secret Access Key:",
		"",
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) != len(expected) {
		t.Errorf("missing lines \n%+v\n\nvs\n\n%+v", expected, lines)
	} else {
		for i, e := range expected {
			if strings.TrimSpace(lines[i]) != e {
				t.Errorf("expected %s got %s", e, lines[i])
			}
		}
	}

	if err := db.QueryRow("SELECT count(*) FROM v_profiles;").Scan(&count); err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Errorf("expected count 1 got %d", count)
	}
}

func TestListNoProfiles(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	profiles, err := state.ListProfiles()
	if err != nil {
		t.Error(err)
	}

	if len(profiles) != 0 {
		t.Errorf("expected to profiles got %d profiles", len(profiles))
	}

	expecetd := `+------+-------------+-------+--------+
| NAME | DESCRIPTION | CLOUD | ACTIVE |
+------+-------------+-------+--------+
+------+-------------+-------+--------+`
	if profiles.String() != expecetd {
		t.Errorf("expecetd:\n%s\ngot:\n%s", expecetd, profiles.String())
	}
}

func TestListAwsCloud(t *testing.T) {
	_, cleaner := database()
	defer cleaner()
	err := data.CreateProfile(&tinycloud.Profile{
		Name:        "test-name",
		Description: "test-des",
		Settings: &tinycloud.CloudSettings{
			Aws: &tinycloud.AwsSettings{
				AwsRegion:           "reg-1",
				AwsAccessKeyId:      "xxxx-xxx-xxx",
				AwsSeacretAccessKey: "yyy-yyy-yy",
			},
		},
	})

	if err != nil {
		t.Error(err)
	}

	profiles, err := state.ListProfiles()
	if err != nil {
		t.Error(err)
	}

	expecetd := `+-----------+-------------+-------+--------+
| NAME      | DESCRIPTION | CLOUD | ACTIVE |
+-----------+-------------+-------+--------+
| test-name | test-des    | aws   | x      |
+-----------+-------------+-------+--------+`
	if profiles.String() != expecetd {
		t.Errorf("expecetd:\n%s\ngot:\n%s", expecetd, profiles.String())
	}
}

func TestListMultipleProfiles(t *testing.T) {
	_, cleaner := database()
	defer cleaner()
	err := data.CreateProfile(&tinycloud.Profile{
		Name:        "test-name",
		Description: "test-des",
		Settings: &tinycloud.CloudSettings{
			Aws: &tinycloud.AwsSettings{
				AwsRegion:           "reg-1",
				AwsAccessKeyId:      "xxxx-xxx-xxx",
				AwsSeacretAccessKey: "yyy-yyy-yy",
			},
		},
	})

	if err != nil {
		t.Error(err)
	}

	err = data.CreateProfile(&tinycloud.Profile{
		Name:        "test-name-2",
		Description: "test-des-2",
		Settings: &tinycloud.CloudSettings{
			Aws: &tinycloud.AwsSettings{
				AwsRegion:           "reg0",
				AwsAccessKeyId:      "aaaa-aaaa-aaa",
				AwsSeacretAccessKey: "oooo-oooo-oooo",
			},
		},
	})

	if err != nil {
		t.Error(err)
	}

	profiles, err := state.ListProfiles()
	if err != nil {
		t.Error(err)
	}

	expecetd := `+-------------+-------------+-------+--------+
| NAME        | DESCRIPTION | CLOUD | ACTIVE |
+-------------+-------------+-------+--------+
| test-name-2 | test-des-2  | aws   | x      |
| test-name   | test-des    | aws   |        |
+-------------+-------------+-------+--------+`
	if profiles.String() != expecetd {
		t.Errorf("expecetd:\n%s\ngot:\n%s", expecetd, profiles.String())
	}
}
