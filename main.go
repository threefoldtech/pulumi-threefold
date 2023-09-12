package main

import (
	"log"

	p "github.com/rawdaGastan/pulumi-provider-grid/internal"
)

var version = "v1.0.0"
var providerName = "grid"

func main() {
	if err := p.RunProvider(providerName, version); err != nil {
		log.Fatal(err.Error())
	}
}
