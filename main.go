package main

import (
	"github.com/kucicm/tiny-cloud/cmd"
	"github.com/kucicm/tiny-cloud/pkg/data"
)

func main() {
	_ = data.SetupDatabes("") // ugly
	defer data.CloseDatabes()

	cmd.Execute()
}
