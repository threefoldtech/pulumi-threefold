import os
import pulumi
import pulumi_threefold as threefold

mnemonic = os.environ['MNEMONIC']
network = os.environ['NETWORK']

provider = threefold.Provider("provider", mnemonic=mnemonic, network=network)
scheduler = threefold.Scheduler("scheduler",
    mru=0.25,
    sru=2,
    farm_ids=[1],
    opts = pulumi.ResourceOptions(provider=provider))
deployment = threefold.Deployment("deployment",
    node_id=scheduler.nodes[0],
    name="zdb",
    zdbs=[threefold.ZDBInputArgs(
        name="zdbsTest",
        size=2,
        password="123456",
    )],
    opts = pulumi.ResourceOptions(provider=provider))

pulumi.export("node_deployment_id", deployment.node_deployment_id)
pulumi.export("zdb_endpoint_ip", deployment.zdbs_computed[0].ips[1])
pulumi.export("zdb_endpoint_port", deployment.zdbs_computed[0].port)
pulumi.export("zdb_namespace", deployment.zdbs_computed[0].namespace)
