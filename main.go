package main

import (
	"log"

	p "github.com/rawdaGastan/pulumi-provider-grid/internal"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string
var ProviderName = "pulumi-provider-grid"

func main() {
	if err := p.RunProvider(ProviderName, Version); err != nil {
		log.Println(err)
	}
}
