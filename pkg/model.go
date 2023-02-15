package tinycloud

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Profiles []*Profile

func (ps Profiles) String() string {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Name", "Description", "Cloud", "Active"})
	for _, profile := range []*Profile(ps) {
		tw.AppendRow(profile.toBasicRow())
	}
	tw.SetAutoIndex(false)
	return tw.Render()
}

type Profile struct {
	Name        string
	Description string
	Active      bool
	Settings    *CloudSettings
}

func (p *Profile) toBasicRow() table.Row {
	cloud := p.Settings.ResolveCloudName()
	active := ""
	if p.Active {
		active = "x"
	}
	return table.Row{p.Name, p.Description, cloud, active}
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
	Aws *AwsSettings
	// gcp
}

type AwsSettings struct {
	AwsRegion           string
	AwsAccessKeyId      string
	AwsSeacretAccessKey string
}

func (s *AwsSettings) Valid() error {
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
}

func (s *CloudSettings) Valid() error {
	name := s.ResolveCloudName()
	switch name {
	case "aws":
		return s.Aws.Valid()
	default:
		return fmt.Errorf("unknown cloud '%s'", name)
	}
}

func (s *CloudSettings) ResolveCloudName() string {
	if s.Aws != nil {
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

type RunRequest struct {
	DockerImageId string
	VmType      string
	DataOutPath string
}

type TaskDefinition struct {
	DockerImageId string
	DataOutPath string
	DNSName     string
	SSHKey      []byte
}

type Vm struct {
	Id      string
	SSHKey  []byte
	DNSName string
}
