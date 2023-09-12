package provider

import (
	"fmt"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
)

func parseToFqdnGatewayState(fqdnGateway workloads.GatewayFQDNProxy) FqdnGatewayState {

	stateArgs := FqdnGatewayArgs{
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
	for key, value := range fqdnGateway.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(key)] = int64(value)
	}

	state := FqdnGatewayState{
		FqdnGatewayArgs:  stateArgs,
		ContractID:       int64(fqdnGateway.ContractID),
		NodeDeploymentID: nodeDeploymentID,
	}

	return state

}

func parseToWorkloadFqdnGateway(fqdnGatewayArgs FqdnGatewayArgs) workloads.GatewayFQDNProxy {

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

func parseToFqdnGatewayComputed(fqdnGatewayState FqdnGatewayState) workloads.GatewayFQDNProxy {

	nodeDeploymentID := make(map[uint32]uint64)

	for key, value := range fqdnGatewayState.NodeDeploymentID {
		k, err := strconv.ParseUint(key, 10, 32)
		if err != nil {
			continue
		}
		nodeDeploymentID[uint32(k)] = uint64(value)
	}

	return workloads.GatewayFQDNProxy{
		NodeDeploymentID: nodeDeploymentID,
		ContractID:       uint64(fqdnGatewayState.ContractID),
	}

}

func updateFqdnGatewayFromState(fqdnGateway *workloads.GatewayFQDNProxy, state FqdnGatewayState) error {

	nodeDeploymentID := make(map[uint32]uint64)

	for key, value := range state.NodeDeploymentID {
		k, err := strconv.ParseUint(key, 10, 32)
		if err != nil {
			continue
		}
		nodeDeploymentID[uint32(k)] = uint64(value)
	}

	fqdnGateway.ContractID = uint64(state.ContractID)
	fqdnGateway.NodeDeploymentID = nodeDeploymentID

	return nil

}
