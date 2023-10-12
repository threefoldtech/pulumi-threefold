package main

import (
	"log"

	p "github.com/threefoldtech/pulumi-threefold/provider"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string
var providerName = "threefold"

func main() {
	if err := p.RunProvider(providerName, Version); err != nil {
		log.Fatal(err.Error())
	}
}
