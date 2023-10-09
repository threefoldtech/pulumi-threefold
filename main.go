package main

import (
	"log"

	p "github.com/threefoldtech/pulumi-provider-grid/internal"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string
var providerName = "grid"

func main() {
	if err := p.RunProvider(providerName, Version); err != nil {
		log.Fatal(err.Error())
	}
}
