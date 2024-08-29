package provider

import (
	"context"
	"fmt"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// Kubernetes controlling struct
type Kubernetes struct{}

// KubernetesArgs is defining what arguments it accepts
type KubernetesArgs struct {
	Master       K8sNodeInput   `pulumi:"master"`
	Workers      []K8sNodeInput `pulumi:"workers"`
	Token        string         `pulumi:"token"`
	NetworkName  string         `pulumi:"network_name"`
	SolutionType string         `pulumi:"solution_type,optional"`
	SSHKey       string         `pulumi:"ssh_key,optional"`
}

// KubernetesState is describing the fields that exist on the created resource.
type KubernetesState struct {
	KubernetesArgs

	MasterComputed   K8sNodeComputed            `pulumi:"master_computed"`
	WorkersComputed  map[string]K8sNodeComputed `pulumi:"workers_computed"`
	NodesIPRange     map[string]string          `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[string]int64           `pulumi:"node_deployment_id"`
}

var _ = (infer.Annotated)((*KubernetesArgs)(nil))

func (k *KubernetesArgs) Annotate(a infer.Annotator) {
	a.SetDefault(&k.SolutionType, fmt.Sprintf("kubernetes/%s", k.Master.Name))
}

// Check validates kubernetes data
func (*Kubernetes) Check(
	ctx context.Context,
	name string, oldInputs,
	newInputs resource.PropertyMap,
) (KubernetesArgs, []p.CheckFailure, error) {
	args, checkFailures, err := infer.DefaultCheck[KubernetesArgs](ctx, newInputs)
	if err != nil {
		return args, checkFailures, err
	}

	// TODO: bypass validation of empty node (will be assigned from scheduler)
	if nodeID, ok := args.Master.Node.(string); ok && len(nodeID) == 0 {
		args.Master.Node = 1
	}

	for i := range args.Workers {
		if nodeID, ok := args.Workers[i].Node.(string); ok && len(nodeID) == 0 {
			args.Workers[i].Node = 1
		}
	}

	Kubernetes, err := parseToK8sCluster(args)
	if err != nil {
		return args, checkFailures, err
	}

	return args, checkFailures, Kubernetes.Validate()
}

// Create creates Kubernetes cluster and deploy it
func (*Kubernetes) Create(
	ctx context.Context,
	id string,
	input KubernetesArgs,
	preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return id, state, nil
	}

	k8sCluster, err := parseToK8sCluster(input)
	if err != nil {
		return id, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return id, state, err
	}

	state = parseToK8sState(k8sCluster)

	return id, state, nil
}

// Update updates the arguments of the Kubernetes resource
func (*Kubernetes) Update(
	ctx context.Context,
	id string,
	oldState KubernetesState,
	input KubernetesArgs,
	preview bool) (KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return state, nil
	}

	k8sCluster, err := parseToK8sCluster(input)
	if err != nil {
		return state, err
	}

	if err := updateK8sFromState(&k8sCluster, oldState); err != nil {
		return state, err
	}

	config := infer.GetConfig[Config](ctx)

	config.TFPluginClient.State.Networks.UpdateNetworkSubnets(k8sCluster.NetworkName, k8sCluster.NodesIPRange)

	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
		return state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return state, err
	}

	state = parseToK8sState(k8sCluster)

	return state, nil
}

// Read get the state of the Kubernetes resource
func (*Kubernetes) Read(ctx context.Context, id string, oldState KubernetesState) (string, KubernetesState, error) {
	k8sCluster, err := parseToK8sCluster(oldState.KubernetesArgs)
	if err != nil {
		return id, oldState, err
	}

	if err := updateK8sFromState(&k8sCluster, oldState); err != nil {
		return id, oldState, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Validate(ctx, &k8sCluster); err != nil {
		return id, oldState, err
	}

	if err := k8sCluster.InvalidateBrokenAttributes(config.TFPluginClient.SubstrateConn); err != nil {
		return id, oldState, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return id, oldState, err
	}

	state := parseToK8sState(k8sCluster)

	return id, state, nil
}

// Delete deletes the Kubernetes resource
func (*Kubernetes) Delete(ctx context.Context, id string, oldState KubernetesState) error {
	k8sCluster, err := parseToK8sCluster(oldState.KubernetesArgs)
	if err != nil {
		return err
	}

	if err := updateK8sFromState(&k8sCluster, oldState); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	return config.TFPluginClient.K8sDeployer.Cancel(ctx, &k8sCluster)
}
