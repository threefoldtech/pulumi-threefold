package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/zos"
)

func TestNetworkParser(t *testing.T) {
	networkInput := NetworkArgs{
		Name:         "network",
		Description:  "description",
		Nodes:        []interface{}{1},
		IPRange:      "1.1.1.1/16",
		AddWGAccess:  false,
		SolutionType: "solution",
	}

	t.Run("parsing input success", func(t *testing.T) {
		network, err := parseToZNet(networkInput, false)
		assert.NoError(t, err)
		assert.Equal(t, network.GetName(), networkInput.Name)
		assert.Equal(t, network.GetDescription(), networkInput.Description)
		assert.Equal(t, network.GetSolutionType(), networkInput.SolutionType)
		assert.Equal(t, network.GetAddWGAccess(), networkInput.AddWGAccess)
	})

	t.Run("parsing input failed: wrong node id type", func(t *testing.T) {
		networkInput.Nodes[0] = "node"
		_, err := parseToZNet(networkInput, false)
		assert.Error(t, err)
		networkInput.Nodes[0] = 1
	})

	t.Run("parsing input failed: wrong ip range", func(t *testing.T) {
		networkInput.IPRange = "wrong"
		_, err := parseToZNet(networkInput, false)
		assert.Error(t, err)
		networkInput.IPRange = "1.1.1.1/16"
	})

	t.Run("parsing and update network success", func(t *testing.T) {
		network, err := parseToZNet(networkInput, false)
		assert.NoError(t, err)

		network.SetNodeDeploymentID(map[uint32]uint64{1: 1})
		ip, err := zos.ParseIPNet("1.1.1.1/16")
		assert.NoError(t, err)
		network.SetNodesIPRange(map[uint32]zos.IPNet{1: ip})
		network.SetExternalIP(&ip)

		state := parseNetworkToState(network)
		assert.Equal(t, network.GetNodeDeploymentID()[1], uint64(state.NodeDeploymentID["1"]))

		err = updateNetworkFromState(network, state)
		assert.NoError(t, err)
	})

	t.Run("parsing and update network failed: wrong node id type", func(t *testing.T) {
		network, err := parseToZNet(networkInput, false)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.NodeDeploymentID = map[string]int64{"": 1}
		err = updateNetworkFromState(network, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update network failed: wrong node id type (ip ranges)", func(t *testing.T) {
		network, err := parseToZNet(networkInput, false)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.NodesIPRange = map[string]string{"": "1.1.1.1/16"}
		err = updateNetworkFromState(network, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update network failed: wrong node ip ranges", func(t *testing.T) {
		network, err := parseToZNet(networkInput, false)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.NodesIPRange = map[string]string{"1": "ip"}
		err = updateNetworkFromState(network, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update network failed: wrong external ip", func(t *testing.T) {
		network, err := parseToZNet(networkInput, false)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.ExternalIP = "-"
		err = updateNetworkFromState(network, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update network failed: wrong external sk", func(t *testing.T) {
		network, err := parseToZNet(networkInput, false)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.ExternalSK = "ips"
		err = updateNetworkFromState(network, state)
		assert.Error(t, err)
	})
}
