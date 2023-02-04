package crud_test

import (
	"bytes"
	"testing"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/crud"
)

func TestPrityPrintProfiles(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	profile := &tinycloud.Profile{Name: "test-profile-1", Description: "test des"}
	err := crud.SaveProfile(profile)
	if err != nil {
		t.Errorf("did not expect error %s", err)
	}

	out := &bytes.Buffer{}
	crud.PrityPrintAllProfiles(out)

	expected := `+----------------+-------------+
| NAME           | DESCRIPTION |
+----------------+-------------+
| test-profile-1 | test des    |
+----------------+-------------+`

	if expected != out.String() {
		t.Errorf("expected:\n%s\n\ngot:\n%s\n", expected, out)
	}

}

func TestCreateNewProfile(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	out := &bytes.Buffer{}
	in := &bytes.Buffer{}

	in.Write([]byte("create-test-1\ntest-2\n"))

	err := crud.CreateNewProfile(in, out)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}

	expected := "Name: \nDescription: \n"
	if out.String() != expected {
		t.Errorf("expected:\n%s\n\ngot:\n%s\n", expected, out)
	}

	out = &bytes.Buffer{}
	crud.PrityPrintAllProfiles(out)
	expected = `+---------------+-------------+
| NAME          | DESCRIPTION |
+---------------+-------------+
| create-test-1 | test-2      |
+---------------+-------------+`
	if expected != out.String() {
		t.Errorf("expected:\n%s\n\ngot:\n%s\n", expected, out)
	}
}
