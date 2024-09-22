package provider

import (
	"fmt"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// TODO: computed

var _ = (infer.Annotated)((*NetworkArgs)(nil))
var _ = (infer.Annotated)((*NetworkState)(nil))

func (n *NetworkArgs) Annotate(a infer.Annotator) {
	a.SetDefault(&n.SolutionType, "Network")

	a.Describe(&n.Name, "The name of the network workload, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported")
	a.Describe(&n.Description, "The description of the network workload, optional with no restrictions")
	a.Describe(&n.SolutionType, "The solution type of the network, displayed as project name in contract metadata")
	a.Describe(&n.Nodes, "The nodes used to deploy the network on, shouldn't be empty")
	a.Describe(&n.IPRange, "The IP range for the network, subnet should be 16")
	a.Describe(&n.Mycelium, "A flag to generate a random mycelium key to support mycelium in the network")
	a.Describe(&n.MyceliumKeys, "A map of nodes as a key and mycelium key for each node, mycelium key length should be 32. Selected nodes must be included in the network's nodes")
	a.Describe(&n.AddWGAccess, "A flag to support wireguard in the network")
}

func (n *NetworkState) Annotate(a infer.Annotator) {
	a.Describe(&n.MyceliumKeys, "A map of nodes as a key and mycelium key for each node, mycelium key length should be 32. Selected nodes must be included in the network's nodes")
	a.Describe(&n.AccessWGConfig, "Generated wireguard configuration for external user access to the network")
	a.Describe(&n.ExternalIP, "Wireguard IP assigned for external user access")
	a.Describe(&n.ExternalSK, "External user private key used in encryption while communicating through Wireguard network")
	a.Describe(&n.PublicNodeID, "Public node id (in case it's added). Used for wireguard access and supporting hidden nodes")
	a.Describe(&n.NodesIPRange, "Computed values of nodes' IP ranges after deployment")
	a.Describe(&n.NodeDeploymentID, "Mapping from each node to its deployment id")
}

var _ = (infer.Annotated)((*DeploymentArgs)(nil))
var _ = (infer.Annotated)((*DeploymentState)(nil))

var _ = (infer.Annotated)((*Mount)(nil))
var _ = (infer.Annotated)((*Zlog)(nil))
var _ = (infer.Annotated)((*Metadata)(nil))
var _ = (infer.Annotated)((*Backend)(nil))
var _ = (infer.Annotated)((*Group)(nil))
var _ = (infer.Annotated)((*Disk)(nil))

// var _ = (infer.Annotated)((*VMInput)(nil))
var _ = (infer.Annotated)((*VMComputed)(nil))
var _ = (infer.Annotated)((*ZDBInput)(nil))
var _ = (infer.Annotated)((*ZDBComputed)(nil))
var _ = (infer.Annotated)((*QSFSInput)(nil))
var _ = (infer.Annotated)((*QSFSComputed)(nil))

func (d *DeploymentArgs) Annotate(a infer.Annotator) {
	a.SetDefault(&d.SolutionType, fmt.Sprintf("vm/%s", d.Name))

	a.Describe(&d.Name, "The name of the deployment, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported")
	a.Describe(&d.NodeID, "The node ID to deploy on, required and should match the requested resources")
	a.Describe(&d.NetworkName, "The name of the network, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported. Network must exist")
	a.Describe(&d.SolutionType, "The name of the solution for created contract to be consistent across threefold tooling (project name in deployment metadata)")
	a.Describe(&d.SolutionProvider, "ID for the deployed solution which allows the creator of the solution to gain a percentage of the rewards")
	a.Describe(&d.Disks, "The disks requested to be included in the deployment")
	a.Describe(&d.ZdbsInputs, "The zdbs requested to be included in the deployment")
	a.Describe(&d.VmsInputs, "The vms requested to be included in the deployment")
	a.Describe(&d.QSFSInputs, "The qsfs instances requested to be included in the deployment")
}

func (d *DeploymentState) Annotate(a infer.Annotator) {
	a.Describe(&d.NodeDeploymentID, "Mapping from each node to its deployment ID")
	a.Describe(&d.ContractID, "The deployment ID")
	a.Describe(&d.IPrange, "IP range of the node for the wireguard network (e.g. 10.1.2.0/24). Has to have a subnet mask of 24")
	a.Describe(&d.ZdbsInputs, "The zdbs output requested to be included in the deployment")
	a.Describe(&d.VmsInputs, "The vms output requested to be included in the deployment")
	a.Describe(&d.QSFSInputs, "The qsfs output instances requested to be included in the deployment")
}

func (v *VMInput) Annotate(a infer.Annotator) {
	a.Describe(&v.Name, "The name of the virtual machine workload, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported")
	a.Describe(&v.NetworkName, "The name of the network, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported. Network must exist")
	a.Describe(&v.Description, "The description of the virtual machine workload, optional with no restrictions")
	a.Describe(&v.NodeID, "The node ID to deploy the virtual machine on, required and should match the requested resources")
	a.Describe(&v.Flist, "The flist to be mounted in the virtual machine, required and should be valid. Example: https://hub.grid.tf/tf-official-apps/base:latest.flist")
	a.Describe(&v.Entrypoint, "The entry point for the flist. Example: /sbin/zinit init")
	a.Describe(&v.FlistChecksum, "The checksum of the flist which should match the checksum of the given flist, optional")
	a.Describe(&v.CPU, "The cpu units needed for the virtual machine. Range in [1: 32]")
	a.Describe(&v.Memory, "The memory capacity for the virtual machine in MB. Min is 250 MB")
	a.Describe(&v.RootfsSize, "The root fs size in GB (type SSD). Can be set as 0 to get the default minimum")
	a.Describe(&v.EnvVars, "The environment variables to be passed to the virtual machine. Example: SSH_KEY")
	a.Describe(&v.GPUs, "A list of gpu IDs to be used in the virtual machine. GPU ID format: <slot>/<vendor>/<device>. Example: 0000:28:00.0/1002/731f")
	a.Describe(&v.Mounts, "A list of mounted disks or volumes")
	a.Describe(&v.Zlogs, "A list of virtual machine loggers")
	a.Describe(&v.Mycelium, "A flag to generate a random mycelium IP seed to support mycelium in the virtual machine")
	a.Describe(&v.MyceliumIPSeed, "The seed used for mycelium IP generated for the virtual machine. It's length should be 6")
	a.Describe(&v.Planetary, "A flag to enable generating a yggdrasil IP for the virtual machine")
	a.Describe(&v.PublicIP, "A flag to enable generating a public IP for the virtual machine, public node is required for it")
	a.Describe(&v.PublicIP6, "A flag to enable generating a public IPv6 for the virtual machine, public node is required for it")
}

func (v *VMComputed) Annotate(a infer.Annotator) {
	a.Describe(&v.MyceliumIPSeed, "The seed used for mycelium IP generated for the virtual machine. It's length should be 6")
	a.Describe(&v.ComputedIP, "The reserved public ipv4 if any")
	a.Describe(&v.ComputedIP6, "The reserved public ipv6 if any")
	a.Describe(&v.PlanetaryIP, "The allocated Yggdrasil IP")
	a.Describe(&v.MyceliumIP, "The allocated mycelium IP")
	a.Describe(&v.ConsoleURL, "The url to access the vm via cloud console on private interface using wireguard")
	a.Describe(&v.IP, "The private wireguard IP of the vm")
}

func (m *Mount) Annotate(a infer.Annotator) {
	a.Describe(&m.Name, "The name of the mounted disk/volume, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported")
	a.Describe(&m.MountPoint, "The mount point of the disk/volume")
}

func (z *Zlog) Annotate(a infer.Annotator) {
	a.Describe(&z.Zmachine, "The name of virtual machine, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported")
	a.Describe(&z.Output, "The output logs URL, should be a valid url")
}

func (z *ZDBInput) Annotate(a infer.Annotator) {
	a.SetDefault(&z.Mode, "user", "")

	a.Describe(&z.Name, "The name of the 0-db workload, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported")
	a.Describe(&z.Description, "The description of the 0-db workload, optional with no restrictions")
	a.Describe(&z.Size, "The 0-db size in GB (type HDD)")
	a.Describe(&z.Password, "The 0-db password")
	a.Describe(&z.Public, "A flag to make 0-db namespace public - readable by anyone")
	a.Describe(&z.Mode, "the enumeration of the modes 0-db can operate in (default user)")
}

func (z *ZDBComputed) Annotate(a infer.Annotator) {
	a.Describe(&z.IPs, "Computed IPs of the ZDB. Two IPs are returned: a public IPv6, and a YggIP, in this order")
	a.Describe(&z.Port, "Port of the ZDB")
	a.Describe(&z.Namespace, "Namespace of the ZDB")
}

func (d *Disk) Annotate(a infer.Annotator) {
	a.Describe(&d.Name, "The name of the disk workload, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported")
	a.Describe(&d.Description, "The description of the disk workload, optional with no restrictions")
	a.Describe(&d.Size, "The disk size in GB (type SSD)")
}

func (q *QSFSInput) Annotate(a infer.Annotator) {
	a.Describe(&q.Name, "The name of the qsfs workload, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported")
	a.Describe(&q.Description, "The description of the qsfs workload, optional with no restrictions")
	a.Describe(&q.Cache, "The size of the fuse mountpoint on the node in MBs (holds qsfs local data before pushing)")
	a.Describe(&q.CompressionAlgorithm, "configuration to use for the compression stage. Currently only snappy is supported")
	a.Describe(&q.EncryptionAlgorithm, "configuration to use for the encryption stage. Currently only AES is supported")
	a.Describe(&q.EncryptionKey, "64 long hex encoded encryption key (e.g. 0000000000000000000000000000000000000000000000000000000000000000)")
	a.Describe(&q.Groups, "The backend groups to write the data to")
	a.Describe(&q.MaxZDBDataDirSize, "Maximum size of the data dir in MiB, if this is set and the sum of the file sizes in the data dir gets higher than this value, the least used, already encoded file will be removed")
	a.Describe(&q.Metadata, "List of ZDB backends configurations")
	a.Describe(&q.MinimalShards, "The minimum amount of shards which are needed to recover the original data")
	a.Describe(&q.ExpectedShards, "The amount of shards which are generated when the data is encoded. Essentially, this is the amount of shards which is needed to be able to recover the data, and some disposable shards which could be lost. The amount of disposable shards can be calculated as expected_shards - minimal_shards")
	a.Describe(&q.RedundantGroups, "The amount of groups which one should be able to loose while still being able to recover the original data")
	a.Describe(&q.RedundantNodes, "The amount of nodes that can be lost in every group while still being able to recover the original data")
}

func (q *QSFSComputed) Annotate(a infer.Annotator) {
	a.Describe(&q.MetricsEndpoint, "Exposed metrics endpoint")

}
func (m *Metadata) Annotate(a infer.Annotator) {
	a.Describe(&m.EncryptionAlgorithm, "configuration to use for the encryption stage. Currently only AES is supported")
	a.Describe(&m.EncryptionKey, "64 long hex encoded encryption key (e.g. 0000000000000000000000000000000000000000000000000000000000000000)")
	a.Describe(&m.Prefix, "Data stored on the remote metadata is prefixed with")
	a.Describe(&m.Type, "configuration for the metadata store to use, currently only ZDB is supported")
	a.Describe(&m.Backends, "List of ZDB backends configurations")
}

func (g *Group) Annotate(a infer.Annotator) {
	a.Describe(&g.Backends, "List of ZDB backends configurations")
}

func (b *Backend) Annotate(a infer.Annotator) {
	a.Describe(&b.Address, "Address of backend ZDB (e.g. [300:a582:c60c:df75:f6da:8a92:d5ed:71ad]:9900 or 60.60.60.60:9900)")
	a.Describe(&b.Namespace, "ZDB namespace")
	a.Describe(&b.Password, "Namespace password")
}

var _ = (infer.Annotated)((*KubernetesArgs)(nil))
var _ = (infer.Annotated)((*KubernetesState)(nil))
var _ = (infer.Annotated)((*K8sNodeInput)(nil))

func (k *KubernetesArgs) Annotate(a infer.Annotator) {
	a.SetDefault(&k.SolutionType, fmt.Sprintf("kubernetes/%s", k.Master.Name))

	a.Describe(&k.Master, "Master holds the configuration of master node in the kubernetes cluster")
	a.Describe(&k.Workers, "Workers is a list holding the workers configuration for the kubernetes cluster")
	a.Describe(&k.NetworkName, "The name of the network, it's required and cannot exceed 50 characters. Only alphanumeric and underscores characters are supported. Network must exist")
	a.Describe(&k.SolutionType, "The solution type of the cluster, displayed as project name in contract metadata")
	a.Describe(&k.SSHKey, "SSH key to access the cluster nodes")
	a.Describe(&k.Token, "The cluster secret token. Each node has to have this token to be part of the cluster. This token should be an alphanumeric non-empty string")
}

func (k *KubernetesState) Annotate(a infer.Annotator) {
	a.Describe(&k.MasterComputed, "The computed fields of the master node")
	a.Describe(&k.WorkersComputed, "List of the computed fields of the worker nodes")
	a.Describe(&k.NodesIPRange, "Computed values of nodes' IP ranges after deployment")
	a.Describe(&k.NodeDeploymentID, "Mapping from each node to its deployment ID")
}

func (k *K8sNodeInput) Annotate(a infer.Annotator) {
	a.Describe(&k.DiskSize, "Data disk size in GBs. Must be between 1GB and 10240GBs (10TBs)")
}

var _ = (infer.Annotated)((*GatewayFQDNArgs)(nil))
var _ = (infer.Annotated)((*GatewayFQDNState)(nil))

func (g *GatewayFQDNArgs) Annotate(a infer.Annotator) {
	a.SetDefault(&g.SolutionType, g.Name)

	a.Describe(&g.Name, "Gateway workload name.  This has to be unique within the deployment. It's required and cannot exceed 50 characters. Must contain only alphanumeric and underscore characters")
	a.Describe(&g.NodeID, "The gateway's node ID")
	a.Describe(&g.FQDN, "The fully qualified domain name of the deployed workload")
	a.Describe(&g.Backends, "The backends of the gateway proxy. must be in the format ip:port if tls_passthrough is set, otherwise the format should be http://ip[:port]")
	a.Describe(&g.NetworkName, "Network name to join, if backend IP is private")
	a.Describe(&g.TLSPassthrough, "TLS passthrough controls the TLS termination, if false, the gateway will terminate the TLS, if True, it will only be terminated by the backend service")
	a.Describe(&g.Description, "The description of the virtual machine workload, optional with no restrictions")
	a.Describe(&g.SolutionType, "The name of the solution for created contract to be consistent across threefold tooling (project name in deployment metadata)")
}

func (g *GatewayFQDNState) Annotate(a infer.Annotator) {
	a.Describe(&g.NodeDeploymentID, "Mapping from each node to its deployment ID")
	a.Describe(&g.ContractID, "The deployment ID")
}

var _ = (infer.Annotated)((*GatewayNameArgs)(nil))
var _ = (infer.Annotated)((*GatewayNameState)(nil))

func (g *GatewayNameArgs) Annotate(a infer.Annotator) {
	a.SetDefault(&g.SolutionType, g.Name)

	a.Describe(&g.Name, "Domain prefix. The fqdn will be <name>.<gateway-domain>. This has to be unique within the deployment. It's required and cannot exceed 50 characters. Must contain only alphanumeric and underscore characters")
	a.Describe(&g.NodeID, "The gateway's node ID")
	a.Describe(&g.Backends, "The backends of the gateway proxy. must be in the format ip:port if tls_passthrough is set, otherwise the format should be http://ip[:port]")
	a.Describe(&g.NetworkName, "Network name to join, if backend IP is private")
	a.Describe(&g.TLSPassthrough, "TLS passthrough controls the TLS termination, if false, the gateway will terminate the TLS, if True, it will only be terminated by the backend service")
	a.Describe(&g.Description, "The description of the virtual machine workload, optional with no restrictions")
	a.Describe(&g.SolutionType, "The name of the solution for created contract to be consistent across threefold tooling (project name in deployment metadata)")
}

func (g *GatewayNameState) Annotate(a infer.Annotator) {
	a.Describe(&g.NodeDeploymentID, "Mapping from each node to its deployment ID")
	a.Describe(&g.ContractID, "The deployment ID")
	a.Describe(&g.NameContractID, "The reserved name contract ID")
	a.Describe(&g.FQDN, "The computed fully qualified domain name of the deployed workload")
}
