package main

import (
	"fmt"
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
			Mru: pulumi.Int(1),
			Sru: pulumi.Int(2),
			Farm_ids: pulumi.IntArray{
				pulumi.Int(1),
			},
		}, pulumi.Provider(tfProvider))
		if err != nil {
			return err
		}
		deployment, err := threefold.NewDeployment(ctx, "deployment", &threefold.DeploymentArgs{
			Node_id: scheduler.Nodes.ApplyT(func(nodes []int) (int, error) {
				return nodes[0], nil
			}).(pulumi.IntOutput),
			Name: pulumi.String("zdb"),
			Zdbs: threefold.ZDBInputArray{
				&threefold.ZDBInputArgs{
					Name:     pulumi.String("zdbsTest"),
					Size:     pulumi.Int(2),
					Password: pulumi.String("123456"),
				},
			},
		}, pulumi.Provider(tfProvider))
		if err != nil {
			return err
		}
		ctx.Export("node_deployment_id", deployment.Node_deployment_id)
		ctx.Export("zdb_endpoint", pulumi.All(deployment.Zdbs_computed, deployment.Zdbs_computed).ApplyT(func(_args []interface{}) (string, error) {
			deploymentZdbs_computed := _args[0].([]threefold.ZDBComputed)
			deploymentZdbs_computed1 := _args[1].([]threefold.ZDBComputed)
			return fmt.Sprintf("[%v]:%v", deploymentZdbs_computed[0].Ips[1], deploymentZdbs_computed1[0].Port), nil
		}).(pulumi.StringOutput))
		ctx.Export("zdb_namespace", deployment.Zdbs_computed.ApplyT(func(zdbs_computed []threefold.ZDBComputed) (*string, error) {
			return &zdbs_computed[0].Namespace, nil
		}).(pulumi.StringPtrOutput))
		return nil
	})
}
