package provider

import (
	"context"
	"fmt"
	"strconv"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
)

// Deployment controlling struct
type Deployment struct{}

// DeploymentArgs is defining what arguments it accepts
type DeploymentArgs struct {
	NodeID           interface{} `pulumi:"node_id"`
	Name             string      `pulumi:"name"`
	NetworkName      string      `pulumi:"network_name,optional"`
	SolutionType     string      `pulumi:"solution_type,optional"`
	SolutionProvider int64       `pulumi:"solution_provider,optional"`
	Disks            []Disk      `pulumi:"disks,optional"`
	ZdbsInputs       []ZDBInput  `pulumi:"zdbs,optional"`
	VmsInputs        []VMInput   `pulumi:"vms,optional"`
	QSFSInputs       []QSFSInput `pulumi:"qsfs,optional"`
}

// DeploymentState is describing the fields that exist on the created resource
type DeploymentState struct {
	DeploymentArgs

	NodeDeploymentID map[string]int64 `pulumi:"node_deployment_id"`
	ContractID       int64            `pulumi:"contract_id"`
	IPrange          string           `pulumi:"ip_range"`
	ZdbsComputed     []ZDBComputed    `pulumi:"zdbs_computed"`
	VmsComputed      []VMComputed     `pulumi:"vms_computed"`
	QsfsComputed     []QSFSComputed   `pulumi:"qsfs_computed"`
}

// Check validates the Deployment
func (*Deployment) Check(
	ctx context.Context,
	name string, oldInputs,
	newInputs property.Map,
) (DeploymentArgs, []p.CheckFailure, error) {
	args, checkFailures, err := infer.DefaultCheck[DeploymentArgs](ctx, newInputs)
	if err != nil {
		return args, checkFailures, err
	}

	// TODO: bypass validation of empty node (will be assigned from scheduler)
	if nodeID, ok := args.NodeID.(string); ok && len(nodeID) == 0 {
		args.NodeID = 1
	}

	for i := range args.VmsInputs {
		if nodeID, ok := args.VmsInputs[i].NodeID.(string); ok && len(nodeID) == 0 {
			args.VmsInputs[i].NodeID = 1
		}
	}

	deployment, err := parseInputToDeployment(args, false)
	if err != nil {
		return args, checkFailures, err
	}

	return args, checkFailures, deployment.Validate()
}

// Create creates a deployment
func (*Deployment) Create(
	ctx context.Context,
	req infer.CreateRequest[DeploymentArgs],
) (infer.CreateResponse[DeploymentState], error) {
	state := DeploymentState{DeploymentArgs: req.Inputs}
	if req.DryRun {
		return infer.CreateResponse[DeploymentState]{ID: req.Name, Output: state}, nil
	}

	config := infer.GetConfig[Config](ctx)

	nodeID, err := strconv.Atoi(fmt.Sprint(req.Inputs.NodeID))
	if err != nil {
		return infer.CreateResponse[DeploymentState]{ID: req.Name, Output: state}, err
	}

	isLight, err := isZosLight(ctx, uint32(nodeID), config.TFPluginClient.NcPool, config.TFPluginClient.SubstrateConn)
	if err != nil {
		return infer.CreateResponse[DeploymentState]{ID: req.Name, Output: state}, err
	}

	deployment, err := parseInputToDeployment(req.Inputs, isLight)
	if err != nil {
		return infer.CreateResponse[DeploymentState]{ID: req.Name, Output: state}, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Deploy(ctx, &deployment); err != nil {
		return infer.CreateResponse[DeploymentState]{ID: req.Name, Output: state}, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return infer.CreateResponse[DeploymentState]{ID: req.Name, Output: state}, err
	}

	state = parseDeploymentToState(deployment)

	return infer.CreateResponse[DeploymentState]{ID: req.Name, Output: state}, nil
}

// Update updates the arguments of the deployment resource
func (*Deployment) Update(
	ctx context.Context,
	req infer.UpdateRequest[DeploymentArgs, DeploymentState],
) (infer.UpdateResponse[DeploymentState], error) {
	state := DeploymentState{DeploymentArgs: req.Inputs}
	if req.DryRun {
		return infer.UpdateResponse[DeploymentState]{Output: state}, nil
	}

	config := infer.GetConfig[Config](ctx)

	nodeID, err := strconv.Atoi(fmt.Sprint(req.Inputs.NodeID))
	if err != nil {
		return infer.UpdateResponse[DeploymentState]{Output: state}, err
	}

	isLight, err := isZosLight(ctx, uint32(nodeID), config.TFPluginClient.NcPool, config.TFPluginClient.SubstrateConn)
	if err != nil {
		return infer.UpdateResponse[DeploymentState]{Output: state}, err
	}

	deployment, err := parseInputToDeployment(req.Inputs, isLight)
	if err != nil {
		return infer.UpdateResponse[DeploymentState]{Output: state}, err
	}

	if err := updateDeploymentFromState(&deployment, req.State, isLight); err != nil {
		return infer.UpdateResponse[DeploymentState]{Output: state}, err
	}

	dl_network := config.TFPluginClient.State.Networks.GetNetwork(deployment.NetworkName)
	dl_network.SetNodeSubnet(deployment.NodeID, deployment.IPrange)

	if err := config.TFPluginClient.DeploymentDeployer.Deploy(ctx, &deployment); err != nil {
		return infer.UpdateResponse[DeploymentState]{Output: state}, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return infer.UpdateResponse[DeploymentState]{Output: state}, err
	}

	state = parseDeploymentToState(deployment)

	return infer.UpdateResponse[DeploymentState]{Output: state}, nil
}

// Read gets the state of the deployment resource
func (*Deployment) Read(ctx context.Context, req infer.ReadRequest[DeploymentArgs, DeploymentState]) (infer.ReadResponse[DeploymentArgs, DeploymentState], error) {
	config := infer.GetConfig[Config](ctx)

	nodeID, err := strconv.Atoi(fmt.Sprint(req.State.NodeID))
	if err != nil {
		return infer.ReadResponse[DeploymentArgs, DeploymentState](req), err
	}

	isLight, err := isZosLight(ctx, uint32(nodeID), config.TFPluginClient.NcPool, config.TFPluginClient.SubstrateConn)
	if err != nil {
		return infer.ReadResponse[DeploymentArgs, DeploymentState](req), err
	}

	deployment, err := parseInputToDeployment(req.State.DeploymentArgs, isLight)
	if err != nil {
		return infer.ReadResponse[DeploymentArgs, DeploymentState](req), err
	}

	if err := updateDeploymentFromState(&deployment, req.State, isLight); err != nil {
		return infer.ReadResponse[DeploymentArgs, DeploymentState](req), err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return infer.ReadResponse[DeploymentArgs, DeploymentState](req), err
	}

	state := parseDeploymentToState(deployment)

	return infer.ReadResponse[DeploymentArgs, DeploymentState]{ID: req.ID, Inputs: req.Inputs, State: state}, nil
}

// Delete deletes a deployment resource
func (*Deployment) Delete(ctx context.Context, req infer.DeleteRequest[DeploymentState]) error {
	config := infer.GetConfig[Config](ctx)

	nodeID, err := strconv.Atoi(fmt.Sprint(req.State.NodeID))
	if err != nil {
		return err
	}

	isLight, err := isZosLight(ctx, uint32(nodeID), config.TFPluginClient.NcPool, config.TFPluginClient.SubstrateConn)
	if err != nil {
		return err
	}

	deployment, err := parseInputToDeployment(req.State.DeploymentArgs, isLight)
	if err != nil {
		return err
	}

	if err := updateDeploymentFromState(&deployment, req.State, isLight); err != nil {
		return err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Cancel(ctx, &deployment); err != nil {
		return err
	}

	return nil
}
