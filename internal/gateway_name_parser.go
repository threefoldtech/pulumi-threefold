package provider

import (
	"fmt"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func parseToGWNameState(gw workloads.GatewayNameProxy) GatewayNameState {

	// parse backends
	backends := make([]string, len(gw.Backends))
	for idx, b := range gw.Backends {
		backends[idx] = string(b)
	}

	// parse NodeDeploymentID
	nodeDeploymentID := make(map[string]int64)
	for nodeID, deploymentID := range gw.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(nodeID)] = int64(deploymentID)
	}

	return GatewayNameState{
		GatewayNameArgs: GatewayNameArgs{
			Name:           gw.Name,
			NodeID:         int32(gw.NodeID),
			Backends:       backends,
			TLSPassthrough: gw.TLSPassthrough,
			Network:        gw.Network,
			Description:    gw.Description,
			SolutionType:   gw.SolutionType,
		},
		NodeDeploymentID: nodeDeploymentID,
		FQDN:             gw.FQDN,
		NameContractID:   int64(gw.NameContractID),
		ContractID:       int64(gw.ContractID),
	}
}

func parseToGWName(gwArgs GatewayNameArgs) workloads.GatewayNameProxy {

	// parse backends
	backends := make([]zos.Backend, len(gwArgs.Backends))
	for idx, b := range gwArgs.Backends {
		backends[idx] = zos.Backend(b)
	}

	return workloads.GatewayNameProxy{
		Name:           gwArgs.Name,
		NodeID:         uint32(gwArgs.NodeID),
		Backends:       backends,
		TLSPassthrough: gwArgs.TLSPassthrough,
		Network:        gwArgs.Network,
		Description:    gwArgs.Description,
		SolutionType:   gwArgs.SolutionType,
	}
}

func updateGWNameFromState(gw *workloads.GatewayNameProxy, state GatewayNameState) error {

	// parse NodeDeploymentID
	nodeDeploymentID := make(map[uint32]uint64)
	for nodeID, deploymentID := range state.NodeDeploymentID {
		node, err := strconv.Atoi(nodeID)
		if err != nil {
			return err
		}

		nodeDeploymentID[uint32(node)] = uint64(deploymentID)
	}

	gw.NodeDeploymentID = nodeDeploymentID
	gw.FQDN = state.FQDN
	gw.NameContractID = uint64(state.NameContractID)
	gw.ContractID = uint64(state.ContractID)

	return nil
}
