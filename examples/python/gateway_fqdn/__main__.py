import os
import pulumi
import pulumi_threefold as threefold

mnemonic = os.environ['MNEMONIC']
network = os.environ['NETWORK']

provider = threefold.Provider("provider", mnemonic=mnemonic, network=network)
scheduler = threefold.Scheduler("scheduler",
    mru=0.25,
    farm_ids=[1],
    ipv4=True,
    free_ips=1,
    opts = pulumi.ResourceOptions(provider=provider))
network = threefold.Network("network",
    name="example",
    description="example network",
    nodes=[scheduler.nodes[0]],
    ip_range="10.1.0.0/16",
    opts = pulumi.ResourceOptions(provider=provider,
        depends_on=[scheduler]))
deployment = threefold.Deployment("deployment",
    node_id=scheduler.nodes[0],
    name="deployment",
    network_name="example",
    vms=[threefold.VMInputArgs(
        name="vm",
        node_id=scheduler.nodes[0],
        flist="https://hub.grid.tf/tf-official-apps/base:latest.flist",
        network_name="example",
        cpu=2,
        memory=256,
        planetary=True,
    )],
    opts = pulumi.ResourceOptions(provider=provider,
        depends_on=[network]))
gateway_fqdn = threefold.GatewayFQDN("gatewayFQDN",
    name="testing",
    node_id=14,
    fqdn="remote.omar.grid.tf",
    backends=[deployment.vms_computed.apply(lambda vms_computed: f"http://[{vms_computed[0].planetary_ip}]:9000")],
    opts = pulumi.ResourceOptions(provider=provider,
        depends_on=[deployment]))
pulumi.export("node_deployment_id", gateway_fqdn.node_deployment_id)
pulumi.export("fqdn", gateway_fqdn.fqdn)
