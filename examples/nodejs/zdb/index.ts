import * as pulumi from "@pulumi/pulumi";
import * as threefold from "@threefold/pulumi";

const provider = new threefold.Provider("provider", {mnemonic: process.env.MNEMONIC, network: process.env.NETWORK});
const scheduler = new threefold.Scheduler("scheduler", {
    mru: 0.25,
    sru: 2,
    farm_ids: [1],
}, {
    provider: provider,
});
const deployment = new threefold.Deployment("deployment", {
    node_id: scheduler.nodes[0],
    name: "zdb",
    zdbs: [{
        name: "zdbsTest",
        size: 2,
        password: "123456",
    }],
}, {
    provider: provider,
});
export const nodeDeploymentId = deployment.node_deployment_id;
export const zdbEndpoint = pulumi.all([deployment.zdbs_computed, deployment.zdbs_computed]).apply(([deploymentZdbs_computed, deploymentZdbs_computed1]) => `[${deploymentZdbs_computed[0].ips?.[1]}]:${deploymentZdbs_computed1[0].port}`);
export const zdbNamespace = deployment.zdbs_computed.apply(zdbs_computed => zdbs_computed[0].namespace);
