package provider

import (
	"fmt"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

func parseToGatewayNameState(gw workloads.GatewayNameProxy) GatewayNameState {

	return GatewayNameState{

	}
}

func parseToGatewayName(gwArgs GatewayNameArgs) workloads.GatewayNameProxy {

	return workloads.GatewayNameProxy{
		// Master:       master,
		// Workers:      workers,
		// Token:        kubernetesArgs.Token,
		// NetworkName:  kubernetesArgs.NetworkName,
		// SolutionType: kubernetesArgs.SolutionType,
		// SSHKey:       kubernetesArgs.SSHKey,
		// NodesIPRange: make(map[uint32]gridtypes.IPNet),
	}
}

func addComputedFieldsFromState(gw *workloads.GatewayNameProxy, state GatewayNameState) error {


	return nil
}
