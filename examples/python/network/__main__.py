import pulumi
import pulumi_threefold as threefold

provider = threefold.Provider("provider", mnemonic=None)
scheduler = threefold.provider.Scheduler("scheduler", farm_ids=[1],
opts=pulumi.ResourceOptions(provider=provider))
network = threefold.provider.Network("network",
    name="testing",
    description="test network",
    nodes=[scheduler.nodes[0]],
    ip_range="10.1.0.0/16",
    opts=pulumi.ResourceOptions(provider=provider,
        depends_on=[scheduler]))
pulumi.export("node_deployment_id", network.node_deployment_id)
pulumi.export("nodes_ip_range", network.nodes_ip_range)
