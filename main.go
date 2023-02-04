package main

import (
	"github.com/kucicm/tiny-cloud/cmd"
	"github.com/kucicm/tiny-cloud/pkg/crud"
)

func main() {
	_ = crud.SetupDatabes("") // ugly
	defer crud.CloseDatabes()

	cmd.Execute()
}
