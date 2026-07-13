package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/threefoldtech/zos_sdk_go/grid-client/zos"
)

// Kubernetes controlling struct
type Kubernetes struct{}

// KubernetesArgs is defining what arguments it accepts
type KubernetesArgs struct {
	Master        K8sNodeInput   `pulumi:"master"`
	Workers       []K8sNodeInput `pulumi:"workers"`
	Token         string         `pulumi:"token"`
	NetworkName   string         `pulumi:"network_name"`
	SolutionType  string         `pulumi:"solution_type,optional"`
	SSHKey        string         `pulumi:"ssh_key,optional"`
	Flist         string         `pulumi:"flist,optional"`
	EntryPoint    string         `pulumi:"entry_point,optional"`
	FlistChecksum string         `pulumi:"flist_checksum,optional"`
}

// KubernetesState is describing the fields that exist on the created resource.
type KubernetesState struct {
	KubernetesArgs

	MasterComputed   VMComputed            `pulumi:"master_computed"`
	WorkersComputed  map[string]VMComputed `pulumi:"workers_computed"`
	NodesIPRange     map[string]string     `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[string]int64      `pulumi:"node_deployment_id"`
}

// Check validates kubernetes data
func (*Kubernetes) Check(
	ctx context.Context,
	name string, oldInputs,
	newInputs property.Map,
) (KubernetesArgs, []p.CheckFailure, error) {
	args, checkFailures, err := infer.DefaultCheck[KubernetesArgs](ctx, newInputs)
	if err != nil {
		return args, checkFailures, err
	}

	// TODO: bypass validation of empty node (will be assigned from scheduler)
	if nodeID, ok := args.Master.NodeID.(string); ok && len(nodeID) == 0 {
		args.Master.NodeID = 1
	}

	for i := range args.Workers {
		if nodeID, ok := args.Workers[i].NodeID.(string); ok && len(nodeID) == 0 {
			args.Workers[i].NodeID = 1
		}
	}

	Kubernetes, err := parseToK8sCluster(args)
	if err != nil {
		return args, checkFailures, err
	}

	// get master and worker flists from the cluster
	if Kubernetes.Flist != "" {
		Kubernetes.Master.Flist = Kubernetes.Flist
		Kubernetes.Master.Entrypoint = Kubernetes.Entrypoint
		for i := range Kubernetes.Workers {
			Kubernetes.Workers[i].Flist = Kubernetes.Flist
			Kubernetes.Workers[i].Entrypoint = Kubernetes.Entrypoint
		}
	}
	return args, checkFailures, Kubernetes.Validate()
}

// Create creates Kubernetes cluster and deploy it
func (*Kubernetes) Create(
	ctx context.Context,
	req infer.CreateRequest[KubernetesArgs],
) (infer.CreateResponse[KubernetesState], error) {
	state := KubernetesState{KubernetesArgs: req.Inputs}
	if req.DryRun {
		return infer.CreateResponse[KubernetesState]{ID: req.Name, Output: state}, nil
	}

	k8sCluster, err := parseToK8sCluster(req.Inputs)
	if err != nil {
		return infer.CreateResponse[KubernetesState]{ID: req.Name, Output: state}, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
		return infer.CreateResponse[KubernetesState]{ID: req.Name, Output: state}, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return infer.CreateResponse[KubernetesState]{ID: req.Name, Output: state}, err
	}

	state = parseToK8sState(k8sCluster)

	return infer.CreateResponse[KubernetesState]{ID: req.Name, Output: state}, nil
}

// Update updates the arguments of the Kubernetes resource
func (*Kubernetes) Update(
	ctx context.Context,
	req infer.UpdateRequest[KubernetesArgs, KubernetesState],
) (infer.UpdateResponse[KubernetesState], error) {
	state := KubernetesState{KubernetesArgs: req.Inputs}
	if req.DryRun {
		return infer.UpdateResponse[KubernetesState]{Output: state}, nil
	}

	k8sCluster, err := parseToK8sCluster(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[KubernetesState]{Output: state}, err
	}

	if err := updateK8sFromState(&k8sCluster, req.State); err != nil {
		return infer.UpdateResponse[KubernetesState]{Output: state}, err
	}

	config := infer.GetConfig[Config](ctx)

	ipRanges := make(map[uint32]zos.IPNet)
	for node, ipRange := range k8sCluster.NodesIPRange {
		ipRanges[node] = zos.IPNet(ipRange)
	}
	config.TFPluginClient.State.Networks.UpdateNetworkSubnets(k8sCluster.NetworkName, ipRanges)

	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
		return infer.UpdateResponse[KubernetesState]{Output: state}, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return infer.UpdateResponse[KubernetesState]{Output: state}, err
	}

	state = parseToK8sState(k8sCluster)

	return infer.UpdateResponse[KubernetesState]{Output: state}, nil
}

// Read get the state of the Kubernetes resource
func (*Kubernetes) Read(ctx context.Context, req infer.ReadRequest[KubernetesArgs, KubernetesState]) (infer.ReadResponse[KubernetesArgs, KubernetesState], error) {
	k8sCluster, err := parseToK8sCluster(req.State.KubernetesArgs)
	if err != nil {
		return infer.ReadResponse[KubernetesArgs, KubernetesState](req), err
	}

	if err := updateK8sFromState(&k8sCluster, req.State); err != nil {
		return infer.ReadResponse[KubernetesArgs, KubernetesState](req), err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Validate(ctx, &k8sCluster); err != nil {
		return infer.ReadResponse[KubernetesArgs, KubernetesState](req), err
	}

	if err := k8sCluster.InvalidateBrokenAttributes(config.TFPluginClient.SubstrateConn); err != nil {
		return infer.ReadResponse[KubernetesArgs, KubernetesState](req), err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return infer.ReadResponse[KubernetesArgs, KubernetesState](req), err
	}

	state := parseToK8sState(k8sCluster)

	return infer.ReadResponse[KubernetesArgs, KubernetesState]{ID: req.ID, Inputs: req.Inputs, State: state}, nil
}

// Delete deletes the Kubernetes resource
func (*Kubernetes) Delete(ctx context.Context, req infer.DeleteRequest[KubernetesState]) error {
	k8sCluster, err := parseToK8sCluster(req.State.KubernetesArgs)
	if err != nil {
		return err
	}

	if err := updateK8sFromState(&k8sCluster, req.State); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	return config.TFPluginClient.K8sDeployer.Cancel(ctx, &k8sCluster)
}
