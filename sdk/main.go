package main

import (
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/threefoldtech/pulumi-threefold/sdk/go/threefold"
	"github.com/threefoldtech/pulumi-threefold/sdk/go/threefold/provider"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		mnemonic := os.Getenv("MNEMONIC")

		_, err := threefold.NewProvider(ctx, "grid provider", &threefold.ProviderArgs{
			Mnemonic: pulumi.String(mnemonic),
		})

		if err != nil {
			return err
		}

		_, err = provider.NewNetwork(ctx, "grid network", &provider.NetworkArgs{
			Description: pulumi.String("example network"),
			Ip_range:    pulumi.String("10.1.0.0/16"),
			Name:        pulumi.String("example"),
			Nodes:       pulumi.Array{pulumi.Int(14)},
		})

		if err != nil {
			return err
		}

		return nil
	})
}
