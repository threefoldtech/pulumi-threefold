package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/zos/pkg/gridtypes"
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
		network, err := parseToZNet(networkInput)
		assert.NoError(t, err)
		assert.Equal(t, network.Name, networkInput.Name)
		assert.Equal(t, network.Nodes[0], uint32(networkInput.Nodes[0].(int)))
		assert.Equal(t, network.Description, networkInput.Description)
		assert.Equal(t, network.SolutionType, networkInput.SolutionType)
		assert.Equal(t, network.AddWGAccess, networkInput.AddWGAccess)
	})

	t.Run("parsing input failed: wrong node id type", func(t *testing.T) {
		networkInput.Nodes[0] = "node"
		_, err := parseToZNet(networkInput)
		assert.Error(t, err)
		networkInput.Nodes[0] = 1
	})

	t.Run("parsing input failed: wrong ip range", func(t *testing.T) {
		networkInput.IPRange = "wrong"
		_, err := parseToZNet(networkInput)
		assert.Error(t, err)
		networkInput.IPRange = "1.1.1.1/16"
	})

	t.Run("parsing and update network success", func(t *testing.T) {
		network, err := parseToZNet(networkInput)
		assert.NoError(t, err)

		network.NodeDeploymentID = map[uint32]uint64{1: 1}
		ip, err := gridtypes.ParseIPNet("1.1.1.1/16")
		assert.NoError(t, err)
		network.NodesIPRange = map[uint32]gridtypes.IPNet{1: ip}
		network.ExternalIP = &ip

		state := parseNetworkToState(network)
		assert.Equal(t, network.NodeDeploymentID[1], uint64(state.NodeDeploymentID["1"]))

		err = updateNetworkFromState(&network, state)
		assert.NoError(t, err)
	})

	t.Run("parsing and update network failed: wrong node id type", func(t *testing.T) {
		network, err := parseToZNet(networkInput)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.NodeDeploymentID = map[string]int64{"": 1}
		err = updateNetworkFromState(&network, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update network failed: wrong node id type (ip ranges)", func(t *testing.T) {
		network, err := parseToZNet(networkInput)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.NodesIPRange = map[string]string{"": "1.1.1.1/16"}
		err = updateNetworkFromState(&network, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update network failed: wrong node ip ranges", func(t *testing.T) {
		network, err := parseToZNet(networkInput)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.NodesIPRange = map[string]string{"1": "ip"}
		err = updateNetworkFromState(&network, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update network failed: wrong external ip", func(t *testing.T) {
		network, err := parseToZNet(networkInput)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.ExternalIP = "-"
		err = updateNetworkFromState(&network, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update network failed: wrong external sk", func(t *testing.T) {
		network, err := parseToZNet(networkInput)
		assert.NoError(t, err)

		state := parseNetworkToState(network)
		state.ExternalSK = "ips"
		err = updateNetworkFromState(&network, state)
		assert.Error(t, err)
	})
}
