package provider

import (
	"context"
	"encoding/hex"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	client "github.com/threefoldtech/tfgrid-sdk-go/grid-client/node"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/subi"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/zos"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func parseNetworkToState(network workloads.Network) NetworkState {
	nodes := []interface{}{}
	for _, nodeID := range network.GetNodes() {
		nodes = append(nodes, nodeID)
	}

	nodesIPRange := make(map[string]string)
	for nodeID, ipRange := range network.GetNodesIPRange() {
		nodesIPRange[fmt.Sprint(nodeID)] = ipRange.String()
	}

	nodeDeploymentID := make(map[string]int64)
	for nodeID, deploymentID := range network.GetNodeDeploymentID() {
		nodeDeploymentID[fmt.Sprint(nodeID)] = int64(deploymentID)
	}

	var myceliumKeys map[string]string
	if len(network.GetMyceliumKeys()) > 0 {
		myceliumKeys = make(map[string]string)
	}
	for nodeID, myceliumKey := range network.GetMyceliumKeys() {
		myceliumKeys[fmt.Sprint(nodeID)] = hex.EncodeToString(myceliumKey)
	}

	stateArgs := NetworkArgs{
		Name:         network.GetName(),
		Description:  network.GetDescription(),
		Nodes:        nodes,
		IPRange:      network.GetIPRange().String(),
		AddWGAccess:  network.GetAddWGAccess(),
		SolutionType: network.GetSolutionType(),
		MyceliumKeys: myceliumKeys,
	}

	state := NetworkState{
		NetworkArgs:      stateArgs,
		AccessWGConfig:   network.GetAccessWGConfig(),
		ExternalSK:       network.GetExternalSK().String(),
		PublicNodeID:     int32(network.GetPublicNodeID()),
		NodesIPRange:     nodesIPRange,
		NodeDeploymentID: nodeDeploymentID,
		MyceliumKeys:     myceliumKeys,
	}

	if network.GetExternalIP() != nil {
		state.ExternalIP = network.GetExternalIP().String()
	}

	return state
}

func parseToZNet(networkArgs NetworkArgs, light bool) (workloads.Network, error) {
	ipRange, err := zos.ParseIPNet(networkArgs.IPRange)
	if err != nil {
		return nil, err
	}

	nodes, err := parseNodes(networkArgs.Nodes)
	if err != nil {
		return nil, err
	}

	myceliumKeys := make(map[uint32][]byte)
	for nodeID, myceliumKey := range networkArgs.MyceliumKeys {
		if len(strings.TrimSpace(nodeID)) == 0 {
			continue
		}

		nodeID, err := strconv.Atoi(nodeID)
		if err != nil {
			return nil, err
		}

		myceliumKey, err := hex.DecodeString(myceliumKey)
		if err != nil {
			return nil, err
		}

		myceliumKeys[uint32(nodeID)] = myceliumKey
	}

	if networkArgs.Mycelium && len(myceliumKeys) == 0 {
		for _, nodeID := range nodes {
			myceliumKey, err := workloads.RandomMyceliumKey()
			if err != nil {
				return nil, err
			}

			myceliumKeys[nodeID] = myceliumKey
		}
	}

	if light {
		return &workloads.ZNetLight{
			Name:         networkArgs.Name,
			Description:  networkArgs.Description,
			Nodes:        nodes,
			IPRange:      ipRange,
			SolutionType: networkArgs.SolutionType,
			MyceliumKeys: myceliumKeys,
		}, nil
	}

	return &workloads.ZNet{
		Name:         networkArgs.Name,
		Description:  networkArgs.Description,
		Nodes:        nodes,
		IPRange:      ipRange,
		AddWGAccess:  networkArgs.AddWGAccess,
		SolutionType: networkArgs.SolutionType,
		MyceliumKeys: myceliumKeys,
	}, nil
}

func updateNetworkFromState(network workloads.Network, state NetworkState) error {
	externalIP, err := zos.ParseIPNet(state.ExternalIP)
	if err != nil {
		return err
	}

	externalSk, err := wgtypes.ParseKey(state.ExternalSK)
	if err != nil {
		return err
	}

	nodesIPRange := make(map[uint32]zos.IPNet)
	for nodeID, ipRange := range state.NodesIPRange {
		ip, err := zos.ParseIPNet(ipRange)
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

	myceliumKeys := make(map[uint32][]byte)
	for nodeID, myceliumKey := range state.MyceliumKeys {
		nodeID, err := strconv.Atoi(fmt.Sprint(nodeID))
		if err != nil {
			return err
		}

		myceliumKey, err := hex.DecodeString(myceliumKey)
		if err != nil {
			return err
		}

		myceliumKeys[uint32(nodeID)] = myceliumKey
	}

	network.SetAccessWGConfig(state.AccessWGConfig)
	network.SetExternalIP(&externalIP)
	network.SetExternalSK(externalSk)
	network.SetPublicNodeID(uint32(state.PublicNodeID))
	network.SetNodesIPRange(nodesIPRange)
	network.SetNodeDeploymentID(nodeDeploymentID)
	network.SetMyceliumKeys(myceliumKeys)

	return nil
}

func isZosLight(ctx context.Context, nodeID uint32, ncPool client.NodeClientGetter, sub subi.SubstrateExt) (bool, error) {
	nodeClient, err := ncPool.GetNodeClient(sub, nodeID)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get node client '%d'", nodeID)
	}

	features, err := nodeClient.SystemGetNodeFeatures(ctx)
	if err != nil {
		return false, errors.Wrapf(err, "failed to get node features '%d'", nodeID)
	}

	return slices.Contains(features, zos.NetworkLightType), nil
}

func isNetworkLight(ctx context.Context, nodeIDs []uint32, ncPool client.NodeClientGetter, sub subi.SubstrateExt) (bool, error) {
	for _, n := range nodeIDs {
		isLight, err := isZosLight(ctx, n, ncPool, sub)
		if err != nil {
			return false, err
		}

		// if a node found that supports version 3 then it is not light
		if !isLight {
			return isLight, nil
		}
	}

	return true, nil
}

func parseNodes(nodeIDs []interface{}) ([]uint32, error) {
	nodes := []uint32{}
	for _, nodeID := range nodeIDs {
		if len(strings.TrimSpace(fmt.Sprint(nodeID))) == 0 {
			continue
		}

		nodeID, err := strconv.Atoi(fmt.Sprint(nodeID))
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, uint32(nodeID))
	}

	return nodes, nil
}
