package main

import (
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/threefoldtech/pulumi-threefold/sdk/go/threefold"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		tfProvider, err := threefold.NewProvider(ctx, "provider", &threefold.ProviderArgs{
			Mnemonic: pulumi.String(os.Getenv("MNEMONIC")),
		})
		if err != nil {
			return err
		}
		scheduler, err := threefold.NewScheduler(ctx, "scheduler", &threefold.SchedulerArgs{
			Farm_ids: pulumi.IntArray{
				pulumi.Int(1),
			},
			Ipv4:     pulumi.Bool(true),
			Free_ips: pulumi.Int(1),
		}, pulumi.Provider(tfProvider))
		if err != nil {
			return err
		}
		gatewayName, err := threefold.NewGatewayName(ctx, "gatewayName", &threefold.GatewayNameArgs{
			Name: pulumi.String("pulumi"),
			Node_id: scheduler.Nodes.ApplyT(func(nodes []int) (int, error) {
				return nodes[0], nil
			}).(pulumi.IntOutput),
			Backends: pulumi.StringArray{
				pulumi.String("http://69.164.223.208"),
			},
		}, pulumi.Provider(tfProvider), pulumi.DependsOn([]pulumi.Resource{
			scheduler,
		}))
		if err != nil {
			return err
		}
		ctx.Export("node_deployment_id", gatewayName.Node_deployment_id)
		ctx.Export("fqdn", gatewayName.Fqdn)
		return nil
	})
}
