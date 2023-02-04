package tinycloud

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Profiles []*Profile

func (ps Profiles) String() string {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Name", "Description", "Cloud"})
	for _, profile := range []*Profile(ps) {
		cloud := profile.Settings.ResolveCloudName()
		tw.AppendRow(table.Row{profile.Name, profile.Description, cloud})
	}
	tw.SetAutoIndex(false)
	return tw.Render()
}

type Profile struct {
	Name        string
	Description string
	Settings    *CloudSettings
}

func (p *Profile) Valid() error {
	if IsStrEmpty(p.Name) {
		return fmt.Errorf("undefined name")
	}

	if p.Settings == nil {
		return fmt.Errorf("cloud settigs not defiend")
	}

	return p.Settings.Valid()
}

type CloudSettings struct {
	// aws
	AwsRegion           string
	AwsAccessKeyId      string
	AwsSeacretAccessKey string

	// gcp
}

func (s *CloudSettings) Valid() error {
	name := s.ResolveCloudName()
	switch name {
	case "aws":
		if IsStrEmpty(s.AwsRegion) {
			return fmt.Errorf("undefined aws region")
		}
		if IsStrEmpty(s.AwsAccessKeyId) {
			return fmt.Errorf("undefined aws access key")
		}
		if IsStrEmpty(s.AwsSeacretAccessKey) {
			return fmt.Errorf("undefined aws seacret access key")
		}
		return nil
	default:
		return fmt.Errorf("unknown cloud '%s'", name)
	}
}

func (s *CloudSettings) ResolveCloudName() string {
	if s.AwsRegion != "" || s.AwsAccessKeyId != "" || s.AwsSeacretAccessKey != "" {
		return "aws"
	}
	return ""
}

var SupportedClouds []string = []string{
	"aws",
	"gcp",
}

func IsStrEmpty(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
