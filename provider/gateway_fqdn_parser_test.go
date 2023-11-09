package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func TestFQDNParser(t *testing.T) {
	fqdnInput := GatewayFQDNArgs{
		NodeID:         1,
		Name:           "name",
		FQDN:           "fqdn",
		Backends:       []zos.Backend{"backend"},
		NetworkName:    "network",
		Description:    "description",
		TLSPassthrough: false,
		SolutionType:   "solution",
	}

	t.Run("parsing input success", func(t *testing.T) {
		fqdn, err := parseToGatewayFQDN(fqdnInput)
		assert.NoError(t, err)
		assert.Equal(t, fqdn.NodeID, uint32(fqdnInput.NodeID.(int)))
		assert.Equal(t, fqdn.Name, fqdnInput.Name)
		assert.Equal(t, fqdn.Network, fqdnInput.NetworkName)
		assert.Equal(t, fqdn.FQDN, fqdnInput.FQDN)
		assert.Equal(t, fqdn.Backends, fqdnInput.Backends)
		assert.Equal(t, fqdn.Description, fqdnInput.Description)
		assert.Equal(t, fqdn.SolutionType, fqdnInput.SolutionType)
		assert.Equal(t, fqdn.TLSPassthrough, fqdnInput.TLSPassthrough)
	})

	t.Run("parsing input failed: wrong node id type", func(t *testing.T) {
		fqdnInput.NodeID = ""
		_, err := parseToGatewayFQDN(fqdnInput)
		assert.Error(t, err)
		fqdnInput.NodeID = 1
	})

	t.Run("parsing and update fqdn success", func(t *testing.T) {
		fqdn, err := parseToGatewayFQDN(fqdnInput)
		assert.NoError(t, err)

		fqdn.NodeDeploymentID = map[uint32]uint64{1: 1}
		state := parseToGatewayFQDNState(fqdn)
		assert.Equal(t, fqdn.NodeDeploymentID[1], uint64(state.NodeDeploymentID["1"]))

		err = updateGatewayFQDNFromState(&fqdn, state)
		assert.NoError(t, err)
	})

	t.Run("parsing and update fqdn failed: wrong node id type", func(t *testing.T) {
		fqdn, err := parseToGatewayFQDN(fqdnInput)
		assert.NoError(t, err)

		state := parseToGatewayFQDNState(fqdn)
		state.NodeDeploymentID = map[string]int64{"": 1}
		err = updateGatewayFQDNFromState(&fqdn, state)
		assert.Error(t, err)
	})
}
