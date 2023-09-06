package provider

import (
	"fmt"
	"strconv"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
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

func updateNetworkFromState(network *workloads.ZNet, state NetworkState) error {
	externalIP, err := gridtypes.ParseIPNet(state.ExternalIP)
	if err != nil {
		return err
	}

	externalSk, err := wgtypes.ParseKey(state.ExternalSK)
	if err != nil {
		return err
	}

	nodesIPRange := make(map[uint32]gridtypes.IPNet)
	for k, v := range state.NodesIPRange {
		ip, err := gridtypes.ParseIPNet(v)
		if err != nil {
			return err
		}

		kInt, err := strconv.Atoi(k)
		if err != nil {
			return err
		}
		nodesIPRange[uint32(kInt)] = ip
	}

	nodeDeploymentID := make(map[uint32]uint64)
	for k, v := range state.NodeDeploymentID {
		kInt, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		nodeDeploymentID[uint32(kInt)] = uint64(v)
	}

	network.AccessWGConfig = state.AccessWGConfig
	network.ExternalIP = &externalIP
	network.ExternalSK = externalSk
	network.PublicNodeID = uint32(state.PublicNodeID)
	network.NodesIPRange = nodesIPRange
	network.NodeDeploymentID = nodeDeploymentID

	return nil
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
