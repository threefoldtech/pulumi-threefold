package main

import (
	"log"

	p "github.com/rawdaGastan/pulumi-provider-grid/internal"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string
var ProviderName = "grid"

func main() {
	if err := p.RunProvider(ProviderName, "v1.0.0"); err != nil {
		log.Println(err)
	}
}
// go build github.com/rawdaGastan/pulumi-provider-grid/ 
// pulumi plugin install resource pulumi-provider-grid v1.0.0 -f pulumi-provider-grid --reinstall 
//  (cd consumer && pulumi up)