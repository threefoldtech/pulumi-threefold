package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func TestNameGatewayParser(t *testing.T) {
	nameGatewayInput := GatewayNameArgs{
		Name:           "name",
		NodeID:         1,
		Backends:       []string{"backend"},
		TLSPassthrough: false,
		NetworkName:    "network",
		Description:    "description",
		SolutionType:   "solution",
	}

	t.Run("parsing input success", func(t *testing.T) {
		nameGateway, err := parseToGWName(nameGatewayInput)
		assert.NoError(t, err)
		assert.Equal(t, nameGateway.NodeID, uint32(nameGatewayInput.NodeID.(int)))
		assert.Equal(t, nameGateway.Name, nameGatewayInput.Name)
		assert.Equal(t, nameGateway.Backends[0], zos.Backend(nameGatewayInput.Backends[0]))
		assert.Equal(t, nameGateway.Network, nameGatewayInput.NetworkName)
		assert.Equal(t, nameGateway.Description, nameGatewayInput.Description)
		assert.Equal(t, nameGateway.SolutionType, nameGatewayInput.SolutionType)
		assert.Equal(t, nameGateway.TLSPassthrough, nameGatewayInput.TLSPassthrough)
	})

	t.Run("parsing input failed: wrong node id type", func(t *testing.T) {
		nameGatewayInput.NodeID = ""
		_, err := parseToGWName(nameGatewayInput)
		assert.Error(t, err)
		nameGatewayInput.NodeID = 1
	})

	t.Run("parsing and update nameGateway success", func(t *testing.T) {
		nameGateway, err := parseToGWName(nameGatewayInput)
		assert.NoError(t, err)

		nameGateway.NodeDeploymentID = map[uint32]uint64{1: 1}
		state := parseToGWNameState(nameGateway)
		assert.Equal(t, nameGateway.NodeDeploymentID[1], uint64(state.NodeDeploymentID["1"]))

		err = updateGWNameFromState(&nameGateway, state)
		assert.NoError(t, err)
	})

	t.Run("parsing and update nameGateway failed: wrong node id type", func(t *testing.T) {
		nameGateway, err := parseToGWName(nameGatewayInput)
		assert.NoError(t, err)

		state := parseToGWNameState(nameGateway)
		state.NodeDeploymentID = map[string]int64{"": 1}
		err = updateGWNameFromState(&nameGateway, state)
		assert.Error(t, err)
	})
}
