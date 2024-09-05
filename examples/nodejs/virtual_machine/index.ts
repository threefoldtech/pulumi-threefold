import * as threefold from "@threefold/pulumi";

const provider = new threefold.Provider("provider", {mnemonic: process.env.MNEMONIC, network: process.env.NETWORK});
const scheduler = new threefold.Scheduler("scheduler", {
    mru: 0.25,
    sru: 2,
    farm_ids: [1],
}, {
    provider: provider,
});
const network = new threefold.Network("network", {
    name: "test",
    description: "test network",
    nodes: [scheduler.nodes[0]],
    ip_range: "10.1.0.0/16",
    mycelium: true,
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
        entrypoint: "/sbin/zinit init",
        network_name: "test",
        cpu: 2,
        memory: 256,
        planetary: true,
        mycelium: true,
        mounts: [{
            disk_name: "data",
            mount_point: "/app",
        }],
        env_vars: {
            SSH_KEY: "",
        },
    }],
    disks: [{
        name: "data",
        size: 2,
    }],
}, {
    provider: provider,
    dependsOn: [network],
});
export const nodeDeploymentId = deployment.node_deployment_id;
export const planetaryIp = deployment.vms_computed.apply(vms_computed => vms_computed[0].planetary_ip);
export const myceliumIp = deployment.vms_computed.apply(vms_computed => vms_computed[0].mycelium_ip);
