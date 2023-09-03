package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Network controlling struct
type Network struct{}

// NetworkArgs is defining what arguments it accepts
type NetworkArgs struct {
	Name        string          `pulumi:"name"`
	Description string          `pulumi:"description"`
	Nodes       []uint32        `pulumi:"nodes"`
	IPRange     gridtypes.IPNet `pulumi:"ip_range"`
	AddWGAccess bool            `pulumi:"add_wg_access"`
}

// NetworkState is describing the fields that exist on the created resource
type NetworkState struct {
	NetworkArgs

	SolutionType     string                     `pulumi:"solution_type"`
	AccessWGConfig   string                     `pulumi:"access_wg_config"`
	ExternalIP       *gridtypes.IPNet           `pulumi:"external_ip"`
	ExternalSK       wgtypes.Key                `pulumi:"external_sk"`
	PublicNodeID     uint32                     `pulumi:"public_node_id"`
	NodesIPRange     map[uint32]gridtypes.IPNet `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[uint32]uint64          `pulumi:"node_deployment_id"`
	WGPort           map[uint32]int             `pulumi:"wg_port"`
	Keys             map[uint32]wgtypes.Key     `pulumi:"keys"`
}

// Create creates network and deploy it
func (Network) Create(ctx p.Context, tfPluginClient deployer.TFPluginClient, input NetworkArgs, preview bool) (NetworkState, error) {
	state := NetworkState{NetworkArgs: input}
	if preview {
		return state, nil
	}

	network := workloads.ZNet{
		Name:        state.Name,
		Description: state.Description,
		Nodes:       state.Nodes,
		IPRange:     state.IPRange,
		AddWGAccess: state.AddWGAccess,
	}

	if err := tfPluginClient.NetworkDeployer.Deploy(ctx, &network); err != nil {
		return state, err
	}

	// update state
	state.SolutionType = network.SolutionType
	state.AccessWGConfig = network.AccessWGConfig
	state.ExternalIP = network.ExternalIP
	state.ExternalSK = network.ExternalSK
	state.PublicNodeID = network.PublicNodeID
	state.NodesIPRange = network.NodesIPRange
	state.NodeDeploymentID = network.NodeDeploymentID
	state.WGPort = network.WGPort
	state.Keys = network.Keys

	return state, nil
}
