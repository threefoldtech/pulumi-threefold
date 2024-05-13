---
title: threefold
meta_desc: Provides an overview of the threefold grid Provider for Pulumi.
layout: package
---

The Threefold Resource Provider for the [threefold grid](https://threefold.io) lets you manage your infrastructure using Pulumi.

## Example

{{< chooser language "go,yaml" >}}

{{% choosable language go %}}

```go
package main

import (
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
      Farm_ids: pulumi.IntArray{
        pulumi.Int(1),
      },
    }, pulumi.Provider(tfProvider))
    if err != nil {
      return err
    }

    network, err := provider.NewNetwork(ctx, "network", &provider.NetworkArgs{
      Name:        pulumi.String("testing"),
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

    ctx.Export("node_deployment_id", network.Node_deployment_id)
    ctx.Export("nodes_ip_range", network.Nodes_ip_range)
    return nil
  })
}
```

{{% /choosable %}}

{{% choosable language yaml %}}

```yml
name: pulumi-threefold
runtime: yaml

resources:
  provider:
    type: pulumi:providers:threefold
    properties:
      mnemonic:

  scheduler:
    type: threefold:provider:Scheduler
    options:
      provider: ${provider}
    properties:
      farm_ids: [1]

  network:
    type: threefold:provider:Network
    options:
      provider: ${provider}
      dependsOn:
        - ${scheduler}
    properties:
      name: testing
      description: test network
      nodes:
        - ${scheduler.nodes[0]}
      ip_range: 10.1.0.0/16

outputs:
  node_deployment_id: ${network.node_deployment_id}
  nodes_ip_range: ${network.nodes_ip_range}
```

{{% /choosable %}}

{{< /chooser >}}
