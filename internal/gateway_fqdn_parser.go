package provider

import (
	"fmt"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
)

func parseToGatewayFQDNState(fqdnGateway workloads.GatewayFQDNProxy) GatewayFQDNState {
	stateArgs := GatewayFQDNArgs{
		NodeID:         int32(fqdnGateway.NodeID),
		Name:           fqdnGateway.Name,
		Description:    fqdnGateway.Description,
		SolutionType:   fqdnGateway.SolutionType,
		NetworkName:    fqdnGateway.Network,
		TLSPassthrough: fqdnGateway.TLSPassthrough,
		FQDN:           fqdnGateway.FQDN,
		Backends:       fqdnGateway.Backends,
	}

	nodeDeploymentID := make(map[string]int64)
	for nodeID, deploymentID := range fqdnGateway.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(nodeID)] = int64(deploymentID)
	}

	state := GatewayFQDNState{
		GatewayFQDNArgs:  stateArgs,
		ContractID:       int64(fqdnGateway.ContractID),
		NodeDeploymentID: nodeDeploymentID,
	}

	return state
}

func parseToGatewayFQDN(fqdnGatewayArgs GatewayFQDNArgs) workloads.GatewayFQDNProxy {
	return workloads.GatewayFQDNProxy{
		NodeID:         uint32(fqdnGatewayArgs.NodeID),
		Name:           fqdnGatewayArgs.Name,
		SolutionType:   fqdnGatewayArgs.SolutionType,
		Network:        fqdnGatewayArgs.NetworkName,
		FQDN:           fqdnGatewayArgs.FQDN,
		Backends:       fqdnGatewayArgs.Backends,
		TLSPassthrough: fqdnGatewayArgs.TLSPassthrough,
		Description:    fqdnGatewayArgs.Description,
	}
}

func updateGatewayFQDNFromState(fqdnGateway *workloads.GatewayFQDNProxy, state GatewayFQDNState) error {
	nodeDeploymentID := make(map[uint32]uint64)

	for nodeID, deploymentID := range state.NodeDeploymentID {
		node, err := strconv.ParseUint(nodeID, 10, 32)
		if err != nil {
			return err
		}
		nodeDeploymentID[uint32(node)] = uint64(deploymentID)
	}

	fqdnGateway.ContractID = uint64(state.ContractID)
	fqdnGateway.NodeDeploymentID = nodeDeploymentID

	return nil
}
