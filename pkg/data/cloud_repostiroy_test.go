package data_test

import (
	"strings"
	"testing"

	tinycloud "github.com/kucicm/tiny-cloud/pkg"
	"github.com/kucicm/tiny-cloud/pkg/data"
)

func TestGetNewRunId(t *testing.T) {
	_, cleaner := database()
	defer cleaner()

	profile := &tinycloud.Profile{
		Name: "test-name-1",
		Settings: &tinycloud.CloudSettings{
			AwsRegion:           "region-1",
			AwsAccessKeyId:      "access-key-1",
			AwsSeacretAccessKey: "seacret-acc-key-1",
		},
	}

	if err := data.CreateProfile(profile); err != nil {
		t.Error(err)
	}

	unique := make(map[string]bool)
	for i := 0; i < 10; i++ {
		id, err := data.GetNewRunId(profile.Name)
		if err != nil {
			t.Error(err)
		}

		if _, exist := unique[id]; exist {
			t.Errorf("id already exists %s", id)
		} else {
			unique[id] = true
		}

		if !strings.HasPrefix(id, profile.Name) {
			t.Errorf("id %s dose not have profix %s", id, profile.Name)
		}
	}

}
