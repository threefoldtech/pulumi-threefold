package provider

import (
	"context"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
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
	newInputs property.Map,
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

	network, err := parseToZNet(args, false)
	if err != nil {
		return args, checkFailures, err
	}

	return args, checkFailures, network.Validate()
}

// Create creates network and deploy it
func (*Network) Create(ctx context.Context, req infer.CreateRequest[NetworkArgs]) (infer.CreateResponse[NetworkState], error) {
	state := NetworkState{NetworkArgs: req.Inputs}
	if req.DryRun {
		return infer.CreateResponse[NetworkState]{ID: req.Name, Output: state}, nil
	}

	config := infer.GetConfig[Config](ctx)

	nodes, err := parseNodes(req.Inputs.Nodes)
	if err != nil {
		return infer.CreateResponse[NetworkState]{ID: req.Name, Output: state}, nil
	}

	light, err := isNetworkLight(ctx, nodes, config.TFPluginClient.NcPool, config.TFPluginClient.SubstrateConn)
	if err != nil {
		return infer.CreateResponse[NetworkState]{ID: req.Name, Output: state}, err
	}

	network, err := parseToZNet(req.Inputs, light)
	if err != nil {
		return infer.CreateResponse[NetworkState]{ID: req.Name, Output: state}, err
	}

	if err := config.TFPluginClient.NetworkDeployer.Deploy(ctx, network); err != nil {
		return infer.CreateResponse[NetworkState]{ID: req.Name, Output: state}, err
	}

	state = parseNetworkToState(network)

	return infer.CreateResponse[NetworkState]{ID: req.Name, Output: state}, nil
}

// Update updates the arguments of the network resource
func (*Network) Update(
	ctx context.Context,
	req infer.UpdateRequest[NetworkArgs, NetworkState],
) (infer.UpdateResponse[NetworkState], error) {
	state := NetworkState{NetworkArgs: req.Inputs}
	if req.DryRun {
		return infer.UpdateResponse[NetworkState]{Output: state}, nil
	}

	config := infer.GetConfig[Config](ctx)

	nodes, err := parseNodes(req.Inputs.Nodes)
	if err != nil {
		return infer.UpdateResponse[NetworkState]{Output: state}, nil
	}

	light, err := isNetworkLight(ctx, nodes, config.TFPluginClient.NcPool, config.TFPluginClient.SubstrateConn)
	if err != nil {
		return infer.UpdateResponse[NetworkState]{Output: state}, err
	}

	network, err := parseToZNet(req.Inputs, light)
	if err != nil {
		return infer.UpdateResponse[NetworkState]{Output: state}, err
	}

	if err := updateNetworkFromState(network, req.State); err != nil {
		return infer.UpdateResponse[NetworkState]{Output: state}, err
	}

	if err := config.TFPluginClient.NetworkDeployer.Deploy(ctx, network); err != nil {
		return infer.UpdateResponse[NetworkState]{Output: state}, err
	}

	state = parseNetworkToState(network)

	return infer.UpdateResponse[NetworkState]{Output: state}, nil
}

// Read get the state of the network resource
func (*Network) Read(ctx context.Context, req infer.ReadRequest[NetworkArgs, NetworkState]) (infer.ReadResponse[NetworkArgs, NetworkState], error) {
	config := infer.GetConfig[Config](ctx)

	nodes, err := parseNodes(req.State.Nodes)
	if err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState](req), nil
	}

	light, err := isNetworkLight(ctx, nodes, config.TFPluginClient.NcPool, config.TFPluginClient.SubstrateConn)
	if err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState](req), err
	}

	network, err := parseToZNet(req.State.NetworkArgs, light)
	if err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState](req), err
	}

	if err := updateNetworkFromState(network, req.State); err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState](req), err
	}

	config.TFPluginClient.State.Networks.UpdateNetworkSubnets(network.GetName(), network.GetNodesIPRange())

	if err := network.InvalidateBrokenAttributes(config.TFPluginClient.SubstrateConn, config.TFPluginClient.NcPool); err != nil {
		return infer.ReadResponse[NetworkArgs, NetworkState](req), err
	}

	state := parseNetworkToState(network)

	return infer.ReadResponse[NetworkArgs, NetworkState]{ID: req.ID, Inputs: req.Inputs, State: state}, nil
}

// Delete deletes the network resource
func (*Network) Delete(ctx context.Context, req infer.DeleteRequest[NetworkState]) error {
	config := infer.GetConfig[Config](ctx)

	nodes, err := parseNodes(req.State.Nodes)
	if err != nil {
		return err
	}

	light, err := isNetworkLight(ctx, nodes, config.TFPluginClient.NcPool, config.TFPluginClient.SubstrateConn)
	if err != nil {
		return err
	}

	network, err := parseToZNet(req.State.NetworkArgs, light)
	if err != nil {
		return err
	}

	if err := updateNetworkFromState(network, req.State); err != nil {
		return err
	}

	return config.TFPluginClient.NetworkDeployer.Cancel(ctx, network)
}
