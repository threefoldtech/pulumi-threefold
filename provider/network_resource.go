package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
)

// Network controlling struct
type Network struct{}

// NetworkArgs is defining what arguments it accepts
type NetworkArgs struct {
	Name         string            `pulumi:"name"`
	Description  string            `pulumi:"description"`
	Nodes        []interface{}     `pulumi:"nodes"`
	IPRange      string            `pulumi:"ip_range"`
	AddWGAccess  bool              `pulumi:"add_wg_access,optional"`
	SolutionType string            `pulumi:"solution_type,optional"`
	MyceliumKeys map[string]string `pulumi:"mycelium_keys,optional"`
	Mycelium     bool              `pulumi:"mycelium,optional"`
}

// NetworkState is describing the fields that exist on the created resource.
type NetworkState struct {
	NetworkArgs

	MyceliumKeys     map[string]string `pulumi:"mycelium_keys,optional"`
	AccessWGConfig   string            `pulumi:"access_wg_config"`
	ExternalIP       string            `pulumi:"external_ip"`
	ExternalSK       string            `pulumi:"external_sk"`
	PublicNodeID     int32             `pulumi:"public_node_id"`
	NodesIPRange     map[string]string `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[string]int64  `pulumi:"node_deployment_id"`
}

// Check validates the network
func (*Network) Check(
	ctx context.Context,
	name string, oldInputs,
	newInputs resource.PropertyMap,
) (NetworkArgs, []p.CheckFailure, error) {
	args, checkFailures, err := infer.DefaultCheck[NetworkArgs](ctx, newInputs)
	if err != nil {
		return args, checkFailures, err
	}

	// TODO: bypass validation of empty nodes (will be assigned from scheduler)
	for i, node := range args.Nodes {
		if nodeID, ok := node.(string); ok && len(nodeID) == 0 {
			args.Nodes[i] = i + 1
		}
	}

	network, err := parseToZNet(args)
	if err != nil {
		return args, checkFailures, err
	}

	return args, checkFailures, network.Validate()
}

// Create creates network and deploy it
func (*Network) Create(ctx context.Context, id string, input NetworkArgs, preview bool) (string, NetworkState, error) {
	state := NetworkState{NetworkArgs: input}
	if preview {
		return id, state, nil
	}

	network, err := parseToZNet(input)
	if err != nil {
		return id, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.NetworkDeployer.Deploy(ctx, &network); err != nil {
		return id, state, err
	}

	state = parseNetworkToState(network)

	return id, state, nil
}

// Update updates the arguments of the network resource
func (*Network) Update(
	ctx context.Context,
	id string,
	oldState NetworkState,
	input NetworkArgs,
	preview bool) (NetworkState, error) {
	state := NetworkState{NetworkArgs: input}
	if preview {
		return state, nil
	}

	network, err := parseToZNet(input)
	if err != nil {
		return state, err
	}

	if err := updateNetworkFromState(&network, oldState); err != nil {
		return state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.NetworkDeployer.Deploy(ctx, &network); err != nil {
		return state, err
	}

	state = parseNetworkToState(network)

	return state, nil
}

// Read get the state of the network resource
func (*Network) Read(ctx context.Context, id string, oldState NetworkState) (string, NetworkState, error) {
	network, err := parseToZNet(oldState.NetworkArgs)
	if err != nil {
		return id, oldState, err
	}

	if err := updateNetworkFromState(&network, oldState); err != nil {
		return id, oldState, err
	}

	config := infer.GetConfig[Config](ctx)

	config.TFPluginClient.State.Networks.UpdateNetworkSubnets(network.Name, network.NodesIPRange)

	if err := config.TFPluginClient.NetworkDeployer.InvalidateBrokenAttributes(&network); err != nil {
		return id, oldState, err
	}

	if err = config.TFPluginClient.NetworkDeployer.ReadNodesConfig(ctx, &network); err != nil {
		return id, oldState, err
	}

	state := parseNetworkToState(network)

	return id, state, nil
}

// Delete deletes the network resource
func (*Network) Delete(ctx context.Context, id string, oldState NetworkState) error {
	network, err := parseToZNet(oldState.NetworkArgs)
	if err != nil {
		return err
	}

	if err := updateNetworkFromState(&network, oldState); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	return config.TFPluginClient.NetworkDeployer.Cancel(ctx, &network)
}
