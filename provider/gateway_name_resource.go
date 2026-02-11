package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
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
	newInputs property.Map,
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
	req infer.CreateRequest[GatewayNameArgs],
) (infer.CreateResponse[GatewayNameState], error) {
	state := GatewayNameState{GatewayNameArgs: req.Inputs}
	if req.DryRun {
		return infer.CreateResponse[GatewayNameState]{ID: req.Name, Output: state}, nil
	}

	gw, err := parseToGWName(req.Inputs)
	if err != nil {
		return infer.CreateResponse[GatewayNameState]{ID: req.Name, Output: state}, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayNameDeployer.Deploy(ctx, &gw); err != nil {
		return infer.CreateResponse[GatewayNameState]{ID: req.Name, Output: state}, err
	}

	if err := config.TFPluginClient.GatewayNameDeployer.Sync(ctx, &gw); err != nil {
		return infer.CreateResponse[GatewayNameState]{ID: req.Name, Output: state}, err
	}

	state = parseToGWNameState(gw)

	return infer.CreateResponse[GatewayNameState]{ID: req.Name, Output: state}, nil
}

// Update updates the GatewayName resource
func (*GatewayName) Update(
	ctx context.Context,
	req infer.UpdateRequest[GatewayNameArgs, GatewayNameState],
) (infer.UpdateResponse[GatewayNameState], error) {
	state := GatewayNameState{GatewayNameArgs: req.Inputs}
	if req.DryRun {
		return infer.UpdateResponse[GatewayNameState]{Output: state}, nil
	}

	gw, err := parseToGWName(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[GatewayNameState]{Output: state}, err
	}

	if err := updateGWNameFromState(&gw, req.State); err != nil {
		return infer.UpdateResponse[GatewayNameState]{Output: state}, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayNameDeployer.Deploy(ctx, &gw); err != nil {
		return infer.UpdateResponse[GatewayNameState]{Output: state}, err
	}

	if err := config.TFPluginClient.GatewayNameDeployer.Sync(ctx, &gw); err != nil {
		return infer.UpdateResponse[GatewayNameState]{Output: state}, err
	}

	state = parseToGWNameState(gw)

	return infer.UpdateResponse[GatewayNameState]{Output: state}, nil
}

// Read gets the state of the GatewayName resource
func (*GatewayName) Read(ctx context.Context, req infer.ReadRequest[GatewayNameArgs, GatewayNameState]) (infer.ReadResponse[GatewayNameArgs, GatewayNameState], error) {
	gw, err := parseToGWName(req.State.GatewayNameArgs)
	if err != nil {
		return infer.ReadResponse[GatewayNameArgs, GatewayNameState](req), err
	}

	if err := updateGWNameFromState(&gw, req.State); err != nil {
		return infer.ReadResponse[GatewayNameArgs, GatewayNameState](req), err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayNameDeployer.Sync(ctx, &gw); err != nil {
		return infer.ReadResponse[GatewayNameArgs, GatewayNameState](req), err
	}

	state := parseToGWNameState(gw)

	return infer.ReadResponse[GatewayNameArgs, GatewayNameState]{ID: req.ID, Inputs: req.Inputs, State: state}, nil
}

// Delete deletes the GatewayName resource
func (*GatewayName) Delete(ctx context.Context, req infer.DeleteRequest[GatewayNameState]) error {
	gw, err := parseToGWName(req.State.GatewayNameArgs)
	if err != nil {
		return err
	}

	if err := updateGWNameFromState(&gw, req.State); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayNameDeployer.Cancel(ctx, &gw); err != nil {
		return err
	}

	return config.TFPluginClient.GatewayNameDeployer.Sync(ctx, &gw)
}
