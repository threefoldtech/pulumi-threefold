import os
import pulumi
import pulumi_threefold as threefold

mnemonic = os.environ['MNEMONIC']
network = os.environ['NETWORK']

provider = threefold.Provider("provider", mnemonic=mnemonic, network=network)
scheduler = threefold.Scheduler("scheduler",
    mru=6,
    sru=6,
    farm_ids=[1],
    opts = pulumi.ResourceOptions(provider=provider))
network = threefold.Network("network",
    name="test",
    description="test network",
    nodes=[scheduler.nodes[0]],
    ip_range="10.1.0.0/16",
    mycelium=True,
    opts = pulumi.ResourceOptions(provider=provider,
        depends_on=[scheduler]))
kubernetes = threefold.Kubernetes("kubernetes",
    master=threefold.K8sNodeInputArgs(
        name="kubernetes",
        network_name="test",
        node=scheduler.nodes[0],
        disk_size=2,
        planetary=True,
        mycelium=True,
        cpu=2,
        memory=2048,
    ),
    workers=[
        threefold.K8sNodeInputArgs(
            name="worker1",
            network_name="test",
            node=scheduler.nodes[0],
            disk_size=2,
            cpu=2,
            memory=2048,
            mycelium=True,
        ),
        threefold.K8sNodeInputArgs(
            name="worker2",
            network_name="test",
            node=scheduler.nodes[0],
            disk_size=2,
            cpu=2,
            memory=2048,
            mycelium=True,
        ),
    ],
    token="t123456789",
    network_name="test",
    ssh_key=None,
    opts = pulumi.ResourceOptions(provider=provider,
        depends_on=[network]))

pulumi.export("node_deployment_id", kubernetes.node_deployment_id)
pulumi.export("master", kubernetes.master_computed)
