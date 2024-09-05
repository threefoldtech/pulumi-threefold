import * as threefold from "@threefold/pulumi";

const provider = new threefold.Provider("provider", {mnemonic: process.env.MNEMONIC, network: process.env.NETWORK});
const scheduler = new threefold.Scheduler("scheduler", {
    farm_ids: [1],
    ipv4: true,
    free_ips: 1,
}, {
    provider: provider,
});
const gatewayName = new threefold.GatewayName("gatewayName", {
    name: "pulumi",
    node_id: scheduler.nodes[0],
    backends: ["http://69.164.223.208"],
}, {
    provider: provider,
    dependsOn: [scheduler],
});
export const nodeDeploymentId = gatewayName.node_deployment_id;
export const fqdn = gatewayName.fqdn;
