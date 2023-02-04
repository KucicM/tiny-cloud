package state_test

import (
	"bytes"
	"database/sql"
	"os"
	"strings"
	"testing"

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
