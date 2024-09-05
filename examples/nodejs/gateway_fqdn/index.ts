import * as pulumi from "@pulumi/pulumi";
import * as threefold from "@threefold/pulumi";

const provider = new threefold.Provider("provider", {mnemonic: process.env.MNEMONIC, network: process.env.NETWORK});
const scheduler = new threefold.Scheduler("scheduler", {
    mru: 0.25,
    farm_ids: [1],
    ipv4: true,
    free_ips: 1,
}, {
    provider: provider,
});
const network = new threefold.Network("network", {
    name: "test",
    description: "test network",
    nodes: [scheduler.nodes[0]],
    ip_range: "10.1.0.0/16",
}, {
    provider: provider,
    dependsOn: [scheduler],
});
const deployment = new threefold.Deployment("deployment", {
    node_id: scheduler.nodes[0],
    name: "deployment",
    network_name: "test",
    vms: [{
        name: "vm",
        node_id: scheduler.nodes[0],
        flist: "https://hub.grid.tf/tf-official-apps/base:latest.flist",
        network_name: "test",
        cpu: 2,
        memory: 256,
        planetary: true,
    }],
}, {
    provider: provider,
    dependsOn: [network],
});
const gatewayFQDN = new threefold.GatewayFQDN("gatewayFQDN", {
    name: "testing",
    node_id: 14,
    fqdn: "remote.omar.grid.tf",
    backends: [pulumi.interpolate `http://[${deployment.vms_computed[0].planetary_ip}]:9000`],
}, {
    provider: provider,
    dependsOn: [deployment],
});
export const nodeDeploymentId = gatewayFQDN.node_deployment_id;
export const fqdn = gatewayFQDN.fqdn;
