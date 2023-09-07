package provider

import (
	"fmt"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

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
