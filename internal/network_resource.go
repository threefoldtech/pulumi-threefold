package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Import the Pulumi package

// Network controlling struct
type Network struct{}

// NetworkArgs is defining what arguments it accepts
type NetworkArgs struct {
	tfPluginClient deployer.TFPluginClient `pulumi:"tf_plugin_client"`
	Name           string                  `pulumi:"name"`
	Description    string                  `pulumi:"description"`
	Nodes          pulumi.IntArray         `pulumi:"nodes"`
	IPRange        gridtypes.IPNet         `pulumi:"ip_range"`
	AddWGAccess    bool                    `pulumi:"add_wg_access"`
}

// NetworkState is describing the fields that exist on the created resource
type NetworkState struct {
	NetworkArgs

	SolutionType     string                  `pulumi:"solution_type"`
	AccessWGConfig   string                  `pulumi:"access_wg_config"`
	ExternalIP       *gridtypes.IPNet        `pulumi:"external_ip"`
	ExternalSK       wgtypes.Key             `pulumi:"external_sk"`
	PublicNodeID     int                     `pulumi:"public_node_id"`
	NodesIPRange     map[int]gridtypes.IPNet `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[int]int             `pulumi:"node_deployment_id"`
	WGPort           map[int]int             `pulumi:"wg_port"`
	Keys             map[int]wgtypes.Key     `pulumi:"keys"`
}

// Create creates network and deploy it
func (Network) Create(ctx p.Context, name string, input NetworkArgs, preview bool) (string, NetworkState, error) {
	state := NetworkState{NetworkArgs: input}
	if preview {
		return name, state, nil
	}

	network := workloads.ZNet{
		Name:        state.Name,
		Description: state.Description,
		Nodes:       make([]uint32, len(state.Nodes)),
		IPRange:     state.IPRange,
		AddWGAccess: state.AddWGAccess,
	}

	if err := state.tfPluginClient.NetworkDeployer.Deploy(ctx, &network); err != nil {
		return name, state, err
	}

	// update state
	state.SolutionType = network.SolutionType
	state.AccessWGConfig = network.AccessWGConfig
	state.ExternalIP = network.ExternalIP
	state.ExternalSK = network.ExternalSK
	state.PublicNodeID = int(network.PublicNodeID)

	state.NodesIPRange = make(map[int]gridtypes.IPNet)
	for k, v := range network.NodesIPRange {
		state.NodesIPRange[int(k)] = v
	}

	state.NodeDeploymentID = make(map[int]int)
	for k, v := range network.NodeDeploymentID {
		state.NodeDeploymentID[int(k)] = int(v)
	}

	state.WGPort = make(map[int]int)
	for k, v := range network.WGPort {
		state.WGPort[int(k)] = v
	}

	state.Keys = make(map[int]wgtypes.Key)
	for k, v := range network.Keys {
		state.Keys[int(k)] = v
	}

	return name, state, nil
}
