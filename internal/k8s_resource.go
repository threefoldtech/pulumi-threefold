package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

// Kubernetes controlling struct
type Kubernetes struct{}

// KubernetesArgs is defining what arguments it accepts
type KubernetesArgs struct {
	Master           *workloads.K8sNode  `pulumi:"master"`
	Workers          []workloads.K8sNode `pulumi:"workers"`
	Token            string              `pulumi:"token"`
	NetworkName      string              `pulumi:"network_name"`
	SolutionType     string              `pulumi:"solution_type"`
	SSHKey           string              `pulumi:"ssh_key"`
	NodesIPRange     map[int32]string    `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[int32]int64     `pulumi:"node_deployment_id"`
}

type KubernetesState struct {
	KubernetesArgs
}

func parseToK8sCluster(kubernetesArgs KubernetesArgs) (*workloads.K8sCluster, error) {

	// parse NodesIpRange
	nodesIPRange := make(map[uint32]gridtypes.IPNet)
	for k, v := range kubernetesArgs.NodesIPRange {
		ipRange, err := gridtypes.ParseIPNet(v)
		if err != nil {
			return nil, err
		}

		nodesIPRange[uint32(k)] = ipRange
	}

	// parse NodeDeploymentID
	nodeDeploymentId := make(map[uint32]uint64)
	for k, v := range kubernetesArgs.NodeDeploymentID {
		nodeDeploymentId[uint32(k)] = uint64(v)
	}

	k8sCluster := workloads.K8sCluster{
		Master:           kubernetesArgs.Master,
		Workers:          kubernetesArgs.Workers,
		Token:            kubernetesArgs.Token,
		NetworkName:      kubernetesArgs.NetworkName,
		SolutionType:     kubernetesArgs.SolutionType,
		SSHKey:           kubernetesArgs.SSHKey,
		NodesIPRange:     nodesIPRange,
		NodeDeploymentID: nodeDeploymentId,
	}

	return &k8sCluster, nil
}

// Create creates Kubernetes cluster and deploy it
func (*Kubernetes) Create(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster, err := parseToK8sCluster(input)
	if err != nil {
		return name, state, err
	}

	config := infer.GetConfig[Config](ctx)
	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, k8sCluster); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, k8sCluster); err != nil {
		return name, state, err
	}

	return name, state, nil
}

// Update updates the arguments of the Kubernetes resource
func (*Kubernetes) Update(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster, err := parseToK8sCluster(input)
	if err != nil {
		return name, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, k8sCluster); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, k8sCluster); err != nil {
		return name, state, err
	}

	return name, state, nil
}

// Read get the state of the Kubernetes resource
func (*Kubernetes) Read(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster, err := parseToK8sCluster(input)
	if err != nil {
		return name, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Validate(ctx, k8sCluster); err != nil {
		return name, state, err
	}

	if err := k8sCluster.InvalidateBrokenAttributes(config.TFPluginClient.SubstrateConn); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, k8sCluster); err != nil {
		return name, state, err
	}

	return name, state, nil
}

// Delete deletes the Kubernetes resource
func (*Kubernetes) Delete(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster, err := parseToK8sCluster(input)
	if err != nil {
		return name, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Cancel(ctx, k8sCluster); err != nil {
		return name, state, err
	}

	return name, state, nil
}
