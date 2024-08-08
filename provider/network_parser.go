package provider

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func parseNetworkToState(network workloads.ZNet) NetworkState {
	nodes := []interface{}{}
	for _, nodeID := range network.Nodes {
		nodes = append(nodes, nodeID)
	}

	nodesIPRange := make(map[string]string)
	for nodeID, ipRange := range network.NodesIPRange {
		nodesIPRange[fmt.Sprint(nodeID)] = ipRange.String()
	}

	nodeDeploymentID := make(map[string]int64)
	for nodeID, deploymentID := range network.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(nodeID)] = int64(deploymentID)
	}

	var myceliumKeys map[string]string
	if len(network.MyceliumKeys) > 0 {
		myceliumKeys = make(map[string]string)
	}
	for nodeID, myceliumKey := range network.MyceliumKeys {
		myceliumKeys[fmt.Sprint(nodeID)] = hex.EncodeToString(myceliumKey)
	}

	stateArgs := NetworkArgs{
		Name:         network.Name,
		Description:  network.Description,
		Nodes:        nodes,
		IPRange:      network.IPRange.String(),
		AddWGAccess:  network.AddWGAccess,
		SolutionType: network.SolutionType,
		MyceliumKeys: myceliumKeys,
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
	for _, nodeID := range networkArgs.Nodes {
		nodeID, err := strconv.Atoi(fmt.Sprint(nodeID))
		if err != nil {
			return workloads.ZNet{}, err
		}
		nodes = append(nodes, uint32(nodeID))
	}

	myceliumKeys := make(map[uint32][]byte)
	for nodeID, myceliumKey := range networkArgs.MyceliumKeys {
		nodeID, err := strconv.Atoi(fmt.Sprint(nodeID))
		if err != nil {
			return workloads.ZNet{}, err
		}

		myceliumKey, err := hex.DecodeString(myceliumKey)
		if err != nil {
			return workloads.ZNet{}, err
		}

		myceliumKeys[uint32(nodeID)] = myceliumKey
	}

	if networkArgs.Mycelium && len(myceliumKeys) == 0 {
		for _, nodeID := range nodes {
			myceliumKey, err := workloads.RandomMyceliumKey()
			if err != nil {
				return workloads.ZNet{}, err
			}

			myceliumKeys[nodeID] = myceliumKey
		}
	}

	network := workloads.ZNet{
		Name:         networkArgs.Name,
		Description:  networkArgs.Description,
		Nodes:        nodes,
		IPRange:      ipRange,
		AddWGAccess:  networkArgs.AddWGAccess,
		SolutionType: networkArgs.SolutionType,
		MyceliumKeys: myceliumKeys,
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
	for nodeID, ipRange := range state.NodesIPRange {
		ip, err := gridtypes.ParseIPNet(ipRange)
		if err != nil {
			return err
		}

		node, err := strconv.Atoi(nodeID)
		if err != nil {
			return err
		}
		nodesIPRange[uint32(node)] = ip
	}

	nodeDeploymentID := make(map[uint32]uint64)
	for nodeID, deploymentID := range state.NodeDeploymentID {
		node, err := strconv.Atoi(nodeID)
		if err != nil {
			return err
		}

		nodeDeploymentID[uint32(node)] = uint64(deploymentID)
	}

	network.AccessWGConfig = state.AccessWGConfig
	network.ExternalIP = &externalIP
	network.ExternalSK = externalSk
	network.PublicNodeID = uint32(state.PublicNodeID)
	network.NodesIPRange = nodesIPRange
	network.NodeDeploymentID = nodeDeploymentID

	return nil
}
