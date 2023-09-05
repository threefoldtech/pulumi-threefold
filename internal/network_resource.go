package provider

import (
	"fmt"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
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

type NetworkState struct {
	NetworkArgs

	AccessWGConfig   string            `pulumi:"access_wg_config"`
	ExternalIP       string            `pulumi:"external_ip"`
	ExternalSK       string            `pulumi:"external_sk"`
	PublicNodeID     int32             `pulumi:"public_node_id"`
	NodesIPRange     map[string]string `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[string]int64  `pulumi:"node_deployment_id"`
}

func parseToNetworkState(network workloads.ZNet) NetworkState {
	nodes := []int32{}
	for _, v := range network.Nodes {
		nodes = append(nodes, int32(v))
	}

	nodesIPRange := make(map[string]string)
	for k, v := range network.NodesIPRange {
		nodesIPRange[fmt.Sprint(k)] = v.String()
	}

	nodeDeploymentID := make(map[string]int64)
	for k, v := range network.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(k)] = int64(v)
	}

	stateArgs := NetworkArgs{
		Name:         network.Name,
		Description:  network.Description,
		Nodes:        nodes,
		IPRange:      network.IPRange.String(),
		AddWGAccess:  network.AddWGAccess,
		SolutionType: network.SolutionType,
	}

	state := NetworkState{
		NetworkArgs:      stateArgs,
		AccessWGConfig:   network.AccessWGConfig,
		ExternalSK:       network.ExternalSK.String(),
		PublicNodeID:     int32(network.PublicNodeID),
		NodesIPRange:     nodesIPRange,
		NodeDeploymentID: nodeDeploymentID,
	}

	if network.ExternalIP != nil {
		state.ExternalIP = network.ExternalIP.String()
	}

	return state
}

func parseToZNet(networkArgs NetworkArgs) (workloads.ZNet, error) {
	ipRange, err := gridtypes.ParseIPNet(networkArgs.IPRange)
	if err != nil {
		return workloads.ZNet{}, err
	}

	nodes := []uint32{}
	for _, v := range networkArgs.Nodes {
		nodes = append(nodes, uint32(v))
	}

	network := workloads.ZNet{
		Name:         networkArgs.Name,
		Description:  networkArgs.Description,
		Nodes:        nodes,
		IPRange:      ipRange,
		AddWGAccess:  networkArgs.AddWGAccess,
		SolutionType: networkArgs.SolutionType,
	}

	return network, nil
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
func (*Network) Update(ctx p.Context, name string, input NetworkArgs, preview bool) (string, NetworkState, error) {

	state := NetworkState{NetworkArgs: input}
	if preview {
		return name, state, nil
	}
	network, err := parseToZNet(input)
	if err != nil {
		return name, state, err
	}

	// update network
	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.NetworkDeployer.Deploy(ctx, &network); err != nil {
		return name, state, err
	}

	state = parseToNetworkState(network)

	return name, state, nil
}

// ResourceNetworkRead get the state of the network resource
func Read(ctx p.Context, name string, input NetworkArgs, preview bool) (string, NetworkState, error) {

	state := NetworkState{NetworkArgs: input}
	if preview {
		return name, state, nil
	}

	network, err := parseToZNet(input)
	if err != nil {
		return name, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.NetworkDeployer.InvalidateBrokenAttributes(&network); err != nil {
		return name, state, err

	}

	if err = config.TFPluginClient.NetworkDeployer.ReadNodesConfig(ctx, &network); err != nil {
		return name, state, err

	}

	state = parseToNetworkState(network)

	return name, state, nil

}

// Delete deletes the network resource
func Delete(ctx p.Context, name string, input NetworkArgs, preview bool) (string, NetworkState, error) {

	state := NetworkState{NetworkArgs: input}
	if preview {
		return name, state, nil
	}

	network, err := parseToZNet(input)
	if err != nil {
		return name, state, err
	}

	config := infer.GetConfig[Config](ctx)

	if err = config.TFPluginClient.NetworkDeployer.Cancel(ctx, &network); err != nil {
		state = parseToNetworkState(network)
		return name, state, err
	}

	return name, state, nil

}
