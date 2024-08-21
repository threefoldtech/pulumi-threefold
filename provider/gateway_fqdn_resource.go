package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
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

var _ = (infer.Annotated)((*GatewayFQDNArgs)(nil))

func (g *GatewayFQDNArgs) Annotate(a infer.Annotator) {
	a.SetDefault(&g.SolutionType, g.Name)
}

// Check validates fqdn gateway data
func (*GatewayFQDN) Check(
	ctx context.Context,
	name string, oldInputs,
	newInputs resource.PropertyMap,
) (GatewayFQDNArgs, []p.CheckFailure, error) {
	args, checkFailures, err := infer.DefaultCheck[GatewayFQDNArgs](ctx, newInputs)
	if err != nil {
		return args, checkFailures, err
	}

	// TODO: bypass validation of empty node (will be assigned from schedular)
	if nodeID, ok := args.NodeID.(string); ok && len(nodeID) == 0 {
		args.NodeID = 1
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
	id string,
	input GatewayFQDNArgs,
	preview bool) (string, GatewayFQDNState, error) {
	state := GatewayFQDNState{GatewayFQDNArgs: input}
	if preview {
		return id, state, nil
	}

	fqdnGateway, err := parseToGatewayFQDN(input)
	if err != nil {
		return id, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Deploy(ctx, &fqdnGateway); err != nil {
		return id, state, err
	}

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return id, state, err
	}

	state = parseToGatewayFQDNState(fqdnGateway)

	return id, state, nil

}

// Update updates the arguments of a fqdn gateway resource
func (*GatewayFQDN) Update(
	ctx context.Context,
	id string,
	oldState GatewayFQDNState,
	input GatewayFQDNArgs,
	preview bool) (GatewayFQDNState, error) {

	state := GatewayFQDNState{GatewayFQDNArgs: input}
	if preview {
		return state, nil
	}

	fqdnGateway, err := parseToGatewayFQDN(input)
	if err != nil {
		return state, err
	}

	if err := updateGatewayFQDNFromState(&fqdnGateway, oldState); err != nil {
		return state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Deploy(ctx, &fqdnGateway); err != nil {
		return state, err
	}

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return state, err
	}

	state = parseToGatewayFQDNState(fqdnGateway)

	return state, nil

}

// Read gets the state of the fqdn gateway resource
func (*GatewayFQDN) Read(ctx context.Context, id string, oldState GatewayFQDNState) (string, GatewayFQDNState, error) {
	fqdnGateway, err := parseToGatewayFQDN(oldState.GatewayFQDNArgs)
	if err != nil {
		return id, oldState, err
	}

	if err := updateGatewayFQDNFromState(&fqdnGateway, oldState); err != nil {
		return id, oldState, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.GatewayFQDNDeployer.Sync(ctx, &fqdnGateway); err != nil {
		return id, oldState, err
	}

	state := parseToGatewayFQDNState(fqdnGateway)

	return id, state, nil

}

// Delete deletes a fqdn gateway resource
func (*GatewayFQDN) Delete(ctx context.Context, id string, oldState GatewayFQDNState) error {
	fqdnGateway, err := parseToGatewayFQDN(oldState.GatewayFQDNArgs)
	if err != nil {
		return err
	}

	if err := updateGatewayFQDNFromState(&fqdnGateway, oldState); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	return config.TFPluginClient.GatewayFQDNDeployer.Cancel(ctx, &fqdnGateway)
}
