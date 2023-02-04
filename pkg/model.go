package tinycloud

import "github.com/jedib0t/go-pretty/v6/table"

type Profiles []*Profile

func (ps Profiles) String() string {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Name", "Description"})
	for _, profile := range []*Profile(ps) {
		tw.AppendRow(table.Row{profile.Name, profile.Description})
	}
	tw.SetAutoIndex(false)
	return tw.Render()
}

type Profile struct {
	Id          int
	Name        string
	Description string
}
