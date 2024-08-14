import os
import pulumi
import pulumi_threefold as threefold

mnemonic = os.environ['MNEMONIC']
network = os.environ['NETWORK']

provider = threefold.Provider("provider", mnemonic=mnemonic, network=network)
scheduler = threefold.Scheduler("scheduler",
    farm_ids=[1],
    ipv4=True,
    free_ips=1,
    opts = pulumi.ResourceOptions(provider=provider))
gateway_name = threefold.GatewayName("gatewayName",
    name="pulumi",
    node_id=scheduler.nodes[0],
    backends=["http://69.164.223.208"],
    opts = pulumi.ResourceOptions(provider=provider,
        depends_on=[scheduler]))
pulumi.export("node_deployment_id", gateway_name.node_deployment_id)
pulumi.export("fqdn", gateway_name.fqdn)
