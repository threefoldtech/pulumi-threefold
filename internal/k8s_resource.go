package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
)

// Kubernetes controlling struct
type Kubernetes struct{}

// KubernetesArgs is defining what arguments it accepts
type KubernetesArgs struct {
	Master       *workloads.K8sNode  `pulumi:"master"`
	Workers      []workloads.K8sNode `pulumi:"workers"`
	Token        string              `pulumi:"token"`
	NetworkName  string              `pulumi:"network_name"`
	SolutionType string              `pulumi:"solution_type,optional"`
	SSHKey       string              `pulumi:"ssh_key,optional"`
}

type KubernetesState struct {
	KubernetesArgs

	NodesIPRange     map[int32]string `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[int32]int64  `pulumi:"node_deployment_id"`
}

func parseToKubernetesState(k8sCluster workloads.K8sCluster) KubernetesState {

	// parse NodesIpRange
	nodesIPRange := make(map[int32]string)
	for k, v := range k8sCluster.NodesIPRange {
		nodesIPRange[int32(k)] = v.String()
	}

	// parse NodeDeploymentID
	nodeDeploymentID := make(map[int32]int64)
	for k, v := range k8sCluster.NodeDeploymentID {
		nodeDeploymentID[int32(k)] = int64(v)
	}

	return KubernetesState{
		KubernetesArgs: KubernetesArgs{
			Master:       k8sCluster.Master,
			Workers:      k8sCluster.Workers,
			Token:        k8sCluster.Token,
			NetworkName:  k8sCluster.NetworkName,
			SolutionType: k8sCluster.SolutionType,
			SSHKey:       k8sCluster.SSHKey,
		},
		NodesIPRange:     nodesIPRange,
		NodeDeploymentID: nodeDeploymentID,
	}
}

func parseToK8sCluster(kubernetesArgs KubernetesArgs) workloads.K8sCluster {

	return workloads.K8sCluster{
		Master:       kubernetesArgs.Master,
		Workers:      kubernetesArgs.Workers,
		Token:        kubernetesArgs.Token,
		NetworkName:  kubernetesArgs.NetworkName,
		SolutionType: kubernetesArgs.SolutionType,
		SSHKey:       kubernetesArgs.SSHKey,
	}
}

// Create creates Kubernetes cluster and deploy it
func (*Kubernetes) Create(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster := parseToK8sCluster(input)

	config := infer.GetConfig[Config](ctx)
	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}

// Update updates the arguments of the Kubernetes resource
func (*Kubernetes) Update(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster := parseToK8sCluster(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}

// Read get the state of the Kubernetes resource
func (*Kubernetes) Read(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster := parseToK8sCluster(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Validate(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	if err := k8sCluster.InvalidateBrokenAttributes(config.TFPluginClient.SubstrateConn); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}

// Delete deletes the Kubernetes resource
func (*Kubernetes) Delete(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster := parseToK8sCluster(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Cancel(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}
