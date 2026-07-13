package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/threefoldtech/zos_base/pkg/gridtypes/zos"
)

// GatewayFQDN controlling struct
type GatewayFQDN struct{}

// GatewayFQDNArgs is defining what arguments it accepts
type GatewayFQDNArgs struct {
	NodeID         interface{}   `pulumi:"node_id"`
	Name           string        `pulumi:"name"`
	FQDN           string        `pulumi:"fqdn"`
	Backends       []zos.Backend `pulumi:"backends"`
	NetworkName    string        `pulumi:"network_name,optional"`
	Description    string        `pulumi:"description,optional"`
	TLSPassthrough bool          `pulumi:"tls_pass_through,optional"`
	SolutionType   string        `pulumi:"solution_type,optional"`
}

// GatewayFQDNState is describing the fields that exist on the fqdn gateway resource
type GatewayFQDNState struct {
	GatewayFQDNArgs

	ContractID       int64            `pulumi:"contract_id"`
	NodeDeploymentID map[string]int64 `pulumi:"node_deployment_id"`
}

// Check validates fqdn gateway data
func (*GatewayFQDN) Check(
	ctx context.Context,
	name string, oldInputs,
	newInputs property.Map,
) (GatewayFQDNArgs, []p.CheckFailure, error) {
	args, checkFailures, err := infer.DefaultCheck[GatewayFQDNArgs](ctx, newInputs)
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

	gw, err := parseToGatewayFQDN(args)
	if err != nil {
		return args, checkFailures, err
	}

	return args, checkFailures, gw.Validate()
}

// Create creates a fqdn gateway
func (*GatewayFQDN) Create(
	ctx context.Context,
	req infer.CreateRequest[GatewayFQDNArgs],
) (infer.CreateResponse[GatewayFQDNState], error) {
	state := GatewayFQDNState{GatewayFQDNArgs: req.Inputs}
	if req.DryRun {
		return infer.CreateResponse[GatewayFQDNState]{ID: req.Name, Output: state}, nil
	}

	fqdnGateway, err := parseToGatewayFQDN(req.Inputs)
	if err != nil {
		return infer.CreateResponse[GatewayFQDNState]{ID: req.Name, Output: state}, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Deploy(ctx, &fqdnGateway); err != nil {
		return infer.CreateResponse[GatewayFQDNState]{ID: req.Name, Output: state}, err
	}

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return infer.CreateResponse[GatewayFQDNState]{ID: req.Name, Output: state}, err
	}

	state = parseToGatewayFQDNState(fqdnGateway)

	return infer.CreateResponse[GatewayFQDNState]{ID: req.Name, Output: state}, nil
}

// Update updates the arguments of a fqdn gateway resource
func (*GatewayFQDN) Update(
	ctx context.Context,
	req infer.UpdateRequest[GatewayFQDNArgs, GatewayFQDNState],
) (infer.UpdateResponse[GatewayFQDNState], error) {
	state := GatewayFQDNState{GatewayFQDNArgs: req.Inputs}
	if req.DryRun {
		return infer.UpdateResponse[GatewayFQDNState]{Output: state}, nil
	}

	fqdnGateway, err := parseToGatewayFQDN(req.Inputs)
	if err != nil {
		return infer.UpdateResponse[GatewayFQDNState]{Output: state}, err
	}

	if err := updateGatewayFQDNFromState(&fqdnGateway, req.State); err != nil {
		return infer.UpdateResponse[GatewayFQDNState]{Output: state}, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Deploy(ctx, &fqdnGateway); err != nil {
		return infer.UpdateResponse[GatewayFQDNState]{Output: state}, err
	}

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return infer.UpdateResponse[GatewayFQDNState]{Output: state}, err
	}

	state = parseToGatewayFQDNState(fqdnGateway)

	return infer.UpdateResponse[GatewayFQDNState]{Output: state}, nil
}

// Read gets the state of the fqdn gateway resource
func (*GatewayFQDN) Read(ctx context.Context, req infer.ReadRequest[GatewayFQDNArgs, GatewayFQDNState]) (infer.ReadResponse[GatewayFQDNArgs, GatewayFQDNState], error) {
	fqdnGateway, err := parseToGatewayFQDN(req.State.GatewayFQDNArgs)
	if err != nil {
		return infer.ReadResponse[GatewayFQDNArgs, GatewayFQDNState](req), err
	}

	if err := updateGatewayFQDNFromState(&fqdnGateway, req.State); err != nil {
		return infer.ReadResponse[GatewayFQDNArgs, GatewayFQDNState](req), err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return infer.ReadResponse[GatewayFQDNArgs, GatewayFQDNState](req), err
	}

	state := parseToGatewayFQDNState(fqdnGateway)

	return infer.ReadResponse[GatewayFQDNArgs, GatewayFQDNState]{ID: req.ID, Inputs: req.Inputs, State: state}, nil
}

// Delete deletes a fqdn gateway resource
func (*GatewayFQDN) Delete(ctx context.Context, req infer.DeleteRequest[GatewayFQDNState]) error {
	fqdnGateway, err := parseToGatewayFQDN(req.State.GatewayFQDNArgs)
	if err != nil {
		return err
	}

	if err := updateGatewayFQDNFromState(&fqdnGateway, req.State); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	return config.TFPluginClient.GatewayFQDNDeployer.Cancel(ctx, &fqdnGateway)
}
