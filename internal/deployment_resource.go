package provider

import (
	"fmt"
	"strconv"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
)

// Deployment controlling struct
type Deployment struct{}

// DeploymentArgs is defining what arguments it accepts
type DeploymentArgs struct {
	NodeID           int32       `pulumi:"node_id"`
	Name             string      `pulumi:"name"`
	NetworkName      string      `pulumi:"network_name"`
	SolutionType     string      `pulumi:"solution_type,optional"`
	SolutionProvider *int64      `pulumi:"solution_provider,optional"`
	Disks            []Disk      `pulumi:"disks,optional"`
	ZdbsInputs       []ZDBInput  `pulumi:"zdbs_inputs,optional"`
	VmsInputs        []VMInput   `pulumi:"vms_inputs,optional"`
	QSFSInputs       []QSFSInput `pulumi:"qsfs_inputs,optional"`
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

func parseToDeploymentState(deployment workloads.Deployment) DeploymentState {

	var solutionProvider int64

	if deployment.SolutionProvider != nil {
		solutionProvider = int64(*deployment.SolutionProvider)
	}

	stateArgs := DeploymentArgs{

		NodeID:           int32(deployment.NodeID),
		Name:             deployment.Name,
		SolutionType:     deployment.SolutionType,
		SolutionProvider: &solutionProvider,
		NetworkName:      deployment.NetworkName,
		Disks:            convertDisks(deployment.Disks),
		ZdbsInputs:       convertZdbs(deployment.Zdbs),
		VmsInputs:        convertVMs(deployment.Vms),
		QSFSInputs:       convertQSFSs(deployment.QSFS),
	}

	nodeDeploymentID := make(map[string]int64)
	for key, value := range deployment.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(key)] = int64(value)
	}

	state := DeploymentState{

		DeploymentArgs:   stateArgs,
		NodeDeploymentID: nodeDeploymentID,
		ContractID:       int64(deployment.ContractID),
		IPrange:          deployment.IPrange,
		ZdbsComputed:     convertZdbsComputed(deployment.Zdbs),
		VmsComputed:      convertVMsComputed(deployment.Vms),
		QsfsComputed:     convertQSFSsComputed(deployment.QSFS),
	}

	return state
}

func parseToWorkloadDeployment(deploymentArgs DeploymentArgs) workloads.Deployment {

	var solutionProvider *uint64
	if deploymentArgs.SolutionProvider != nil {
		temp := uint64(*deploymentArgs.SolutionProvider)
		solutionProvider = &temp
	}

	return workloads.Deployment{
		NodeID:           uint32(deploymentArgs.NodeID),
		Name:             deploymentArgs.Name,
		SolutionType:     deploymentArgs.SolutionType,
		SolutionProvider: solutionProvider,
		NetworkName:      deploymentArgs.NetworkName,
		Disks:            convertToWorkloadDisks(deploymentArgs.Disks),
		Zdbs:             convertToWorkloadZdbs(deploymentArgs.ZdbsInputs),
		Vms:              convertToWorkloadVMs(deploymentArgs.VmsInputs),
		QSFS:             convertToWorkloadQSFSs(deploymentArgs.QSFSInputs),
	}

}

func parseToComputedDeployment(deploymentState DeploymentState) workloads.Deployment {

	nodeDeploymentID := make(map[uint32]uint64)

	for key, value := range deploymentState.NodeDeploymentID {
		k, err := strconv.ParseUint(key, 10, 32)
		if err != nil {
			continue
		}
		nodeDeploymentID[uint32(k)] = uint64(value)
	}

	return workloads.Deployment{
		NodeDeploymentID: nodeDeploymentID,
		ContractID:       uint64(deploymentState.ContractID),
		IPrange:          deploymentState.IPrange,
	}

}

func updateDeploymentkFromState(deployment *workloads.Deployment, state DeploymentState) error {

	nodeDeploymentID := make(map[uint32]uint64)

	for key, value := range state.NodeDeploymentID {
		k, err := strconv.ParseUint(key, 10, 32)
		if err == nil {
			return err
		}
		if err != nil {
			continue
		}
		nodeDeploymentID[uint32(k)] = uint64(value)
	}

	deployment.NodeDeploymentID = nodeDeploymentID
	deployment.ContractID = uint64(state.ContractID)
	deployment.IPrange = state.IPrange
	deployment.Zdbs = convertToWorkloadZdbsComputed(state.ZdbsComputed)
	deployment.Vms = convertToVMsWorkloadComputed(state.VmsComputed)
	deployment.QSFS = convertToQSFSsWorkloadComputed(state.QsfsComputed)

	return nil
}

// Create creates a deployment
func (*Deployment) Create(ctx p.Context, id string, input DeploymentArgs, preview bool) (string, DeploymentState, error) {

	state := DeploymentState{DeploymentArgs: input}
	if preview {
		return id, state, nil
	}

	deployment := parseToWorkloadDeployment(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Deploy(ctx, &deployment); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return id, state, err
	}

	state = parseToDeploymentState(deployment)

	return id, state, nil
}

// Update updates the arguments of the deployment resource
func (*Deployment) Update(ctx p.Context, id string, input DeploymentArgs, oldState DeploymentState, preview bool) (string, DeploymentState, error) {

	state := DeploymentState{DeploymentArgs: input}
	if preview {
		return id, state, nil
	}

	deployment := parseToWorkloadDeployment(input)

	if err := updateDeploymentkFromState(&deployment, oldState); err != nil {
		return id, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Deploy(ctx, &deployment); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return id, state, err
	}

	state = parseToDeploymentState(deployment)

	return id, state, nil

}

// Read gets the state of the deployment resource
func (*Deployment) Read(ctx p.Context, id string, oldState DeploymentState) (string, DeploymentState, error) {

	deployment := parseToWorkloadDeployment(oldState.DeploymentArgs)

	if err := updateDeploymentkFromState(&deployment, oldState); err != nil {
		return id, oldState, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Sync(ctx, &deployment); err != nil {
		return id, oldState, err
	}

	state := parseToDeploymentState(deployment)

	return id, state, nil

}

// Delete deletes a deployment resource
func (*Deployment) Delete(ctx p.Context, id string, oldState DeploymentState) error {

	deployment := parseToComputedDeployment(oldState)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.DeploymentDeployer.Cancel(ctx, &deployment); err != nil {
		return err
	}

	return nil
}
