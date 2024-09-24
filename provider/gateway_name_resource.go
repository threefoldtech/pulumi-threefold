package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// GatewayName controlling struct
type GatewayName struct{}

// GatewayNameArgs is defining what arguments it accepts
type GatewayNameArgs struct {
	Name           string      `pulumi:"name"`
	NodeID         interface{} `pulumi:"node_id"`
	Backends       []string    `pulumi:"backends"`
	TLSPassthrough bool        `pulumi:"tls_passthrough,optional"`
	NetworkName    string      `pulumi:"network_name,optional"`
	Description    string      `pulumi:"description,optional"`
	SolutionType   string      `pulumi:"solution_type,optional"`
}

// GatewayNameState is describing the fields that exist on the created resource.
type GatewayNameState struct {
	GatewayNameArgs

	NodeDeploymentID map[string]int64 `pulumi:"node_deployment_id"`
	FQDN             string           `pulumi:"fqdn"`
	NameContractID   int64            `pulumi:"name_contract_id"`
	ContractID       int64            `pulumi:"contract_id"`
}

// Check validates name gateway data
func (*GatewayName) Check(
	ctx context.Context,
	name string, oldInputs,
	newInputs resource.PropertyMap,
) (GatewayNameArgs, []p.CheckFailure, error) {
	args, checkFailures, err := infer.DefaultCheck[GatewayNameArgs](ctx, newInputs)
	if err != nil {
		return args, checkFailures, err
	}

	// TODO: bypass validation of empty node (will be assigned from scheduler)
	if nodeID, ok := args.NodeID.(string); ok && len(nodeID) == 0 {
		args.NodeID = 1
	}

	// TODO: bypass validation of empty backend (will be assigned from vm)
	for i, backend := range args.Backends {
		if len(backend) == 0 {
			args.Backends[i] = "http://0.0.0.0"
		}
	}

	gw, err := parseToGWName(args)
	if err != nil {
		return args, checkFailures, err
	}

	return args, checkFailures, gw.Validate()
}

// Create creates GatewayName and deploy it
func (*GatewayName) Create(
	ctx context.Context,
	id string,
	input GatewayNameArgs,
	preview bool) (string, GatewayNameState, error) {
	state := GatewayNameState{GatewayNameArgs: input}
	if preview {
		return id, state, nil
	}

	gw, err := parseToGWName(input)
	if err != nil {
		return id, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayNameDeployer.Deploy(ctx, &gw); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.GatewayNameDeployer.Sync(ctx, &gw); err != nil {
		return id, state, err
	}

	state = parseToGWNameState(gw)

	return id, state, nil
}

// Update updates the GatewayName resource
func (*GatewayName) Update(
	ctx context.Context,
	id string,
	oldState GatewayNameState,
	input GatewayNameArgs,
	preview bool) (GatewayNameState, error) {
	state := GatewayNameState{GatewayNameArgs: input}
	if preview {
		return state, nil
	}

	gw, err := parseToGWName(input)
	if err != nil {
		return state, err
	}

	if err := updateGWNameFromState(&gw, oldState); err != nil {
		return state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayNameDeployer.Deploy(ctx, &gw); err != nil {
		return state, err
	}

	if err := config.TFPluginClient.GatewayNameDeployer.Sync(ctx, &gw); err != nil {
		return state, err
	}

	state = parseToGWNameState(gw)

	return state, nil
}

// Read gets the state of the GatewayName resource
func (*GatewayName) Read(ctx context.Context, id string, oldState GatewayNameState) (string, GatewayNameState, error) {
	gw, err := parseToGWName(oldState.GatewayNameArgs)
	if err != nil {
		return id, oldState, err
	}

	if err := updateGWNameFromState(&gw, oldState); err != nil {
		return id, oldState, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayNameDeployer.Sync(ctx, &gw); err != nil {
		return id, oldState, err
	}

	state := parseToGWNameState(gw)

	return id, state, nil
}

// Delete deletes the GatewayName resource
func (*GatewayName) Delete(ctx context.Context, id string, oldState GatewayNameState) error {
	gw, err := parseToGWName(oldState.GatewayNameArgs)
	if err != nil {
		return err
	}

	if err := updateGWNameFromState(&gw, oldState); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayNameDeployer.Cancel(ctx, &gw); err != nil {
		return err
	}

	return config.TFPluginClient.GatewayNameDeployer.Sync(ctx, &gw)
}
