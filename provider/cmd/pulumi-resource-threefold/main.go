package main

import (
	"log"

	threefold "github.com/threefoldtech/pulumi-threefold/provider"
	"github.com/threefoldtech/pulumi-threefold/provider/pkg/version"
)

var providerName = "threefold"

func main() {
	if err := threefold.RunProvider(providerName, version.Version); err != nil {
		log.Fatal(err.Error())
	}
}
