package main

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/threefoldtech/pulumi-threefold/sdk/go/threefold"
	"github.com/threefoldtech/pulumi-threefold/sdk/go/threefold/provider"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		tfProvider, err := threefold.NewProvider(ctx, "provider", &threefold.ProviderArgs{
			Mnemonic: pulumi.String(os.Getenv("MNEMONIC")),
		})
		if err != nil {
			return err
		}
		scheduler, err := provider.NewScheduler(ctx, "scheduler", &provider.SchedulerArgs{
			Mru: pulumi.Int(1),
			Farm_ids: pulumi.IntArray{
				pulumi.Int(1),
			},
			Ipv4:     pulumi.Bool(true),
			Free_ips: pulumi.Int(1),
		}, pulumi.Provider(tfProvider))
		if err != nil {
			return err
		}
		network, err := provider.NewNetwork(ctx, "network", &provider.NetworkArgs{
			Name:        pulumi.String("test"),
			Description: pulumi.String("test network"),
			Nodes: pulumi.Array{
				scheduler.Nodes.ApplyT(func(nodes []int) (int, error) {
					return nodes[0], nil
				}).(pulumi.IntOutput),
			},
			Ip_range: pulumi.String("10.1.0.0/16"),
		}, pulumi.Provider(tfProvider), pulumi.DependsOn([]pulumi.Resource{
			scheduler,
		}))
		if err != nil {
			return err
		}
		deployment, err := provider.NewDeployment(ctx, "deployment", &provider.DeploymentArgs{
			Node_id: scheduler.Nodes.ApplyT(func(nodes []int) (int, error) {
				return nodes[0], nil
			}).(pulumi.IntOutput),
			Name:         pulumi.String("deployment"),
			Network_name: pulumi.String("test"),
			Vms: provider.VMInputArray{
				&provider.VMInputArgs{
					Name:         pulumi.String("vm"),
					Flist:        pulumi.String("https://hub.grid.tf/tf-official-apps/base:latest.flist"),
					Network_name: pulumi.String("test"),
					Cpu:          pulumi.Int(2),
					Memory:       pulumi.Int(256),
					Planetary:    pulumi.Bool(true),
				},
			},
		}, pulumi.Provider(tfProvider), pulumi.DependsOn([]pulumi.Resource{
			network,
		}))
		if err != nil {
			return err
		}
		gatewayFQDN, err := provider.NewGatewayFQDN(ctx, "gatewayFQDN", &provider.GatewayFQDNArgs{
			Name:    pulumi.String("testing"),
			Node_id: pulumi.Any(14),
			Fqdn:    pulumi.String("remote.omar.grid.tf"),
			Backends: pulumi.StringArray{
				deployment.Vms_computed.ApplyT(func(vms_computed []provider.VMComputed) (string, error) {
					return fmt.Sprintf("http://[%v]:9000", vms_computed[0].Planetary_ip), nil
				}).(pulumi.StringOutput),
			},
		}, pulumi.Provider(tfProvider), pulumi.DependsOn([]pulumi.Resource{
			deployment,
		}))
		if err != nil {
			return err
		}
		ctx.Export("node_deployment_id", gatewayFQDN.Node_deployment_id)
		ctx.Export("fqdn", gatewayFQDN.Fqdn)
		return nil
	})
}
