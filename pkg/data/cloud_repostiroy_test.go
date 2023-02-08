package data_test

import (
	"bytes"
	"strings"
	"testing"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
)

func TestGetNewRunId(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	name, err := createProfile()
	if err != nil {
		t.Error(err)
	}

	unique := make(map[string]bool)
	for i := 0; i < 10; i++ {
		id, err := data.GetNewRunId(name)
		if err != nil {
			t.Error(err)
		}

		if _, exist := unique[id]; exist {
			t.Errorf("id already exists %s", id)
		} else {
			unique[id] = true
		}

		if !strings.HasPrefix(id, name) {
			t.Errorf("id %s dose not have profix %s", id, name)
		}
	}

}

func TestAddPemKey(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	name, err := createProfile()
	if err != nil {
		t.Error(err)
	}

	runId, err := data.GetNewRunId(name)
	if err != nil {
		t.Error(err)
	}

	key := []byte("test-key")
	if err := data.AddPemKey(runId, key); err != nil {
		t.Error(err)
	}

	loadKey, err := data.GetPemKey(runId)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(loadKey, key) {
		t.Errorf("expected %+v got %+v", key, loadKey)
	}

}

func createProfile() (string, error) {
	profile := &tinycloud.Profile{
		Name: "test-name-1",
		Settings: &tinycloud.CloudSettings{
			AwsRegion:           "region-1",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	if err := data.CreateProfile(profile); err != nil {
		return "", err
	}

	return profile.Name, nil
}
