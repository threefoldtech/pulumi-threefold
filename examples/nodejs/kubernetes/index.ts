import * as threefold from "@threefold/pulumi";

const provider = new threefold.Provider("provider", {mnemonic: process.env.MNEMONIC, network: process.env.NETWORK});
const scheduler = new threefold.Scheduler("scheduler", {
    mru: 6,
    sru: 6,
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
const kubernetes = new threefold.Kubernetes("kubernetes", {
    master: {
        name: "kubernetes",
        network_name: "test",
        node_id: scheduler.nodes[0],
        disk_size: 2,
        planetary: true,
        mycelium: true,
        cpu: 2,
        memory: 2048,
    },
    workers: [
        {
            name: "worker1",
            network_name: "test",
            node_id: scheduler.nodes[0],
            disk_size: 2,
            cpu: 2,
            memory: 2048,
        },
        {
            name: "worker2",
            network_name: "test",
            node_id: scheduler.nodes[0],
            disk_size: 2,
            cpu: 2,
            memory: 2048,
        },
    ],
    token: "t123456789",
    network_name: "test",
    ssh_key: undefined,
}, {
    provider: provider,
    dependsOn: [network],
});
export const nodeDeploymentId = kubernetes.node_deployment_id;
export const planetaryIP = kubernetes.master_computed.apply(master_computed => master_computed.planetary_ip);
export const myceliumIP = kubernetes.master_computed.apply(master_computed => master_computed.mycelium_ip);
