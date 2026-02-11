package main

import (
	"context"
	"log"

	threefold "github.com/threefoldtech/pulumi-threefold/provider"
	"github.com/threefoldtech/pulumi-threefold/provider/pkg/version"
)

var providerName = "threefold"

func main() {
	if err := threefold.RunProvider(context.Background(), providerName, version.Version); err != nil {
		log.Fatal(err.Error())
	}
}
