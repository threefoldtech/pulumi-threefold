package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
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

// Annotate sets defaults and descriptions for zdb resource
func (z *ZDBInput) Annotate(a infer.Annotator) {
	a.SetDefault(&z.Mode, "user", "")
}

// Create creates a deployment
func (*Deployment) Create(
	ctx p.Context,
	id string,
	input DeploymentArgs,
	preview bool) (string, DeploymentState, error) {
	state := DeploymentState{DeploymentArgs: input}
	if preview {
		return id, state, nil
	}

	deployment, err := parseInputToDeployment(input)
	if err != nil {
		return id, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Deploy(ctx, &deployment); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return id, state, err
	}

	state = parseDeploymentToState(deployment)

	return id, state, nil
}

// Update updates the arguments of the deployment resource
func (*Deployment) Update(
	ctx p.Context,
	id string,
	oldState DeploymentState,
	input DeploymentArgs,
	preview bool) (string, DeploymentState, error) {
	state := DeploymentState{DeploymentArgs: input}
	if preview {
		return id, state, nil
	}

	deployment, err := parseInputToDeployment(input)
	if err != nil {
		return id, state, err
	}

	if err := updateDeploymentFromState(&deployment, oldState); err != nil {
		return id, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Deploy(ctx, &deployment); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return id, state, err
	}

	state = parseDeploymentToState(deployment)

	return id, state, nil
}

// Read gets the state of the deployment resource
func (*Deployment) Read(ctx p.Context, id string, oldState DeploymentState) (string, DeploymentState, error) {
	deployment, err := parseInputToDeployment(oldState.DeploymentArgs)
	if err != nil {
		return id, oldState, err
	}

	if err := updateDeploymentFromState(&deployment, oldState); err != nil {
		return id, oldState, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return id, oldState, err
	}

	state := parseDeploymentToState(deployment)

	return id, state, nil
}

// Delete deletes a deployment resource
func (*Deployment) Delete(ctx p.Context, id string, oldState DeploymentState) error {
	deployment, err := parseInputToDeployment(oldState.DeploymentArgs)
	if err != nil {
		return err
	}

	if err := updateDeploymentFromState(&deployment, oldState); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Cancel(ctx, &deployment); err != nil {
		return err
	}

	return nil
}
