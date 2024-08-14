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
			Mru: pulumi.Int(6),
			Sru: pulumi.Int(6),
			Farm_ids: pulumi.IntArray{
				pulumi.Int(1),
			},
		}, pulumi.Provider(tfProvider))
		if err != nil {
			return err
		}
		network, err := threefold.NewNetwork(ctx, "network", &threefold.NetworkArgs{
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
		kubernetes, err := threefold.NewKubernetes(ctx, "kubernetes", &threefold.KubernetesArgs{
			Master: &threefold.K8sNodeInputArgs{
				Name: pulumi.String("kubernetes"),
				Node: scheduler.Nodes.ApplyT(func(nodes []int) (int, error) {
					return nodes[0], nil
				}).(pulumi.IntOutput),
				Disk_size: pulumi.Int(2),
				Planetary: pulumi.Bool(true),
				Cpu:       pulumi.Int(2),
				Memory:    pulumi.Int(2048),
			},
			Workers: threefold.K8sNodeInputArray{
				&threefold.K8sNodeInputArgs{
					Name: pulumi.String("worker1"),
					Node: scheduler.Nodes.ApplyT(func(nodes []int) (int, error) {
						return nodes[0], nil
					}).(pulumi.IntOutput),
					Disk_size: pulumi.Int(2),
					Cpu:       pulumi.Int(2),
					Memory:    pulumi.Int(2048),
				},
				&threefold.K8sNodeInputArgs{
					Name: pulumi.String("worker2"),
					Node: scheduler.Nodes.ApplyT(func(nodes []int) (int, error) {
						return nodes[0], nil
					}).(pulumi.IntOutput),
					Disk_size: pulumi.Int(2),
					Cpu:       pulumi.Int(2),
					Memory:    pulumi.Int(2048),
				},
			},
			Token:        pulumi.String("t123456789"),
			Network_name: pulumi.String("test"),
			Ssh_key:      nil,
		}, pulumi.Provider(tfProvider), pulumi.DependsOn([]pulumi.Resource{
			network,
		}))
		if err != nil {
			return err
		}
		ctx.Export("node_deployment_id", kubernetes.Node_deployment_id)
		ctx.Export("planetary_ip", kubernetes.Master_computed.ApplyT(func(master_computed threefold.K8sNodeComputed) (*string, error) {
			return &master_computed.Planetary_ip, nil
		}).(pulumi.StringPtrOutput))
		return nil
	})
}
