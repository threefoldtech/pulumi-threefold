package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// Network controlling struct
type Network struct{}

// NetworkArgs is defining what arguments it accepts
type NetworkArgs struct {
	Name         string  `pulumi:"name"`
	Description  string  `pulumi:"description"`
	Nodes        []int32 `pulumi:"nodes"`
	IPRange      string  `pulumi:"ip_range"`
	AddWGAccess  bool    `pulumi:"add_wg_access,optional"`
	SolutionType string  `pulumi:"solution_type,optional"`
}

// NetworkState is describing the fields that exist on the created resource.
type NetworkState struct {
	NetworkArgs

	AccessWGConfig   string            `pulumi:"access_wg_config"`
	ExternalIP       string            `pulumi:"external_ip"`
	ExternalSK       string            `pulumi:"external_sk"`
	PublicNodeID     int32             `pulumi:"public_node_id"`
	NodesIPRange     map[string]string `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[string]int64  `pulumi:"node_deployment_id"`
}

// Create creates network and deploy it
func (*Network) Create(ctx p.Context, name string, input NetworkArgs, preview bool) (string, NetworkState, error) {

	state := NetworkState{NetworkArgs: input}
	if preview {
		return name, state, nil
	}

	network, err := parseToZNet(input)
	if err != nil {
		return name, state, err
	}

	// deploy network
	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.NetworkDeployer.Deploy(ctx, &network); err != nil {
		return name, state, err
	}

	state = parseToNetworkState(network)

	return name, state, nil
}

// Update updates the arguments of the network resource
func (*Network) Update(ctx p.Context, id string, oldState NetworkState, input NetworkArgs, preview bool) (NetworkState, error) {

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

	// update network
	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.NetworkDeployer.Deploy(ctx, &network); err != nil {
		return state, err
	}

	state = parseToNetworkState(network)

	return state, nil
}

// ResourceNetworkRead get the state of the network resource
func (*Network) Read(ctx p.Context, id string, oldState NetworkState) (string, NetworkState, error) {

	network, err := parseToZNet(oldState.NetworkArgs)
	if err != nil {
		return id, oldState, err
	}
	if err := updateNetworkFromState(&network, oldState); err != nil {
		return id, oldState, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.NetworkDeployer.InvalidateBrokenAttributes(&network); err != nil {
		return id, oldState, err

	}

	if err = config.TFPluginClient.NetworkDeployer.ReadNodesConfig(ctx, &network); err != nil {
		return id, oldState, err

	}

	state := parseToNetworkState(network)

	return id, state, nil
}

// Delete deletes the network resource
func (*Network) Delete(ctx p.Context, id string, oldState NetworkState) error {

	network, err := parseToZNet(oldState.NetworkArgs)
	if err != nil {
		return err
	}
	if err := updateNetworkFromState(&network, oldState); err != nil {
		return err
	}

	config := infer.GetConfig[Config](ctx)

	if err = config.TFPluginClient.NetworkDeployer.Cancel(ctx, &network); err != nil {
		return err
	}

	return nil
}
