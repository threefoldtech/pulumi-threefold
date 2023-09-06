package provider

import (
	"fmt"
	"strconv"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

// FqdnGateway controlling struct
type FqdnGateway struct{}

// FqdnGatewayArgs is defining what arguments it accepts
type FqdnGatewayArgs struct {
	NodeID         int32         `pulumi:"node_id"`
	Name           string        `pulumi:"name"`
	FQDN           string        `pulumi:"fqdn"`
	Backends       []zos.Backend `pulumi:"backends"`
	NetworkName    string        `pulumi:"network_name,optional"`
	Description    string        `pulumi:"description,optional"`
	TLSPassthrough bool          `pulumi:"tls_pass_through,optional"`
	SolutionType   string        `pulumi:"solution_type,optional"`
}

// FqdnGatewayState is describing the fields that exist on the fqdn gateway resource
type FqdnGatewayState struct {
	FqdnGatewayArgs

	ContractID       int64            `pulumi:"contract_id"`
	NodeDeploymentID map[string]int64 `pulumi:"node_deployment_id"`
}

func parseToFqdnGatewayState(fqdnGateway workloads.GatewayFQDNProxy) FqdnGatewayState {

	stateArgs := FqdnGatewayArgs{
		NodeID:         int32(fqdnGateway.NodeID),
		Name:           fqdnGateway.Name,
		Description:    fqdnGateway.Description,
		SolutionType:   fqdnGateway.SolutionType,
		NetworkName:    fqdnGateway.Network,
		TLSPassthrough: fqdnGateway.TLSPassthrough,
		FQDN:           fqdnGateway.FQDN,
		Backends:       fqdnGateway.Backends,
	}

	nodeDeploymentID := make(map[string]int64)
	for key, value := range fqdnGateway.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(key)] = int64(value)
	}

	state := FqdnGatewayState{
		FqdnGatewayArgs:  stateArgs,
		ContractID:       int64(fqdnGateway.ContractID),
		NodeDeploymentID: nodeDeploymentID,
	}

	return state

}

func parseToWorkloadFqdnGateway(fqdnGatewayArgs FqdnGatewayArgs) workloads.GatewayFQDNProxy {

	return workloads.GatewayFQDNProxy{

		NodeID:         uint32(fqdnGatewayArgs.NodeID),
		Name:           fqdnGatewayArgs.Name,
		SolutionType:   fqdnGatewayArgs.SolutionType,
		Network:        fqdnGatewayArgs.NetworkName,
		FQDN:           fqdnGatewayArgs.FQDN,
		Backends:       fqdnGatewayArgs.Backends,
		TLSPassthrough: fqdnGatewayArgs.TLSPassthrough,
		Description:    fqdnGatewayArgs.Description,
	}

}

func parseToFqdnGatewayComputed(fqdnGatewayState FqdnGatewayState) workloads.GatewayFQDNProxy {

	nodeDeploymentID := make(map[uint32]uint64)

	for key, value := range fqdnGatewayState.NodeDeploymentID {
		k, err := strconv.ParseUint(key, 10, 32)
		if err != nil {
			continue
		}
		nodeDeploymentID[uint32(k)] = uint64(value)
	}

	return workloads.GatewayFQDNProxy{
		NodeDeploymentID: nodeDeploymentID,
		ContractID:       uint64(fqdnGatewayState.ContractID),
	}

}

func updateFqdnGatewayFromState(fqdnGateway *workloads.GatewayFQDNProxy, state FqdnGatewayState) error {

	nodeDeploymentID := make(map[uint32]uint64)

	for key, value := range state.NodeDeploymentID {
		k, err := strconv.ParseUint(key, 10, 32)
		if err != nil {
			continue
		}
		nodeDeploymentID[uint32(k)] = uint64(value)
	}

	fqdnGateway.ContractID = uint64(state.ContractID)
	fqdnGateway.NodeDeploymentID = nodeDeploymentID

	return nil

}

// Create creates a fqdn gateway
func (*FqdnGateway) Create(ctx p.Context, id string, input FqdnGatewayArgs, preview bool) (string, FqdnGatewayState, error) {

	state := FqdnGatewayState{FqdnGatewayArgs: input}
	if preview {
		return id, state, nil
	}

	fqdnGateway := parseToWorkloadFqdnGateway(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Deploy(ctx, &fqdnGateway); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return id, state, err
	}

	state = parseToFqdnGatewayState(fqdnGateway)

	return id, state, nil

}

// Update updates the arguments of a fqdn gateway resource
func (*FqdnGateway) Update(ctx p.Context, id string, input FqdnGatewayArgs, oldState FqdnGatewayState, preview bool) (string, FqdnGatewayState, error) {

	state := FqdnGatewayState{FqdnGatewayArgs: input}
	if preview {
		return id, state, nil
	}

	fqdnGateway := parseToWorkloadFqdnGateway(input)

	if err := updateFqdnGatewayFromState(&fqdnGateway, oldState); err != nil {
		return id, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Deploy(ctx, &fqdnGateway); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return id, state, err
	}

	state = parseToFqdnGatewayState(fqdnGateway)

	return id, state, nil

}

// Read gets the state of the fqdn gateway resource
func (*FqdnGateway) Read(ctx p.Context, id string, oldState FqdnGatewayState) (string, FqdnGatewayState, error) {

	fqdnGateway := parseToWorkloadFqdnGateway(oldState.FqdnGatewayArgs)

	if err := updateFqdnGatewayFromState(&fqdnGateway, oldState); err != nil {
		return id, oldState, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return id, oldState, err
	}

	state := parseToFqdnGatewayState(fqdnGateway)

	return id, state, nil

}

// Delete deletes a fqdn gateway resource
func (*FqdnGateway) Delete(ctx p.Context, id string, oldState FqdnGatewayState) error {

	fqdnGateway := parseToFqdnGatewayComputed(oldState)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Cancel(ctx, &fqdnGateway); err != nil {
		return err
	}

	return nil

}
