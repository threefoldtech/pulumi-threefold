package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

func TestK8sParser(t *testing.T) {
	k8sNodeInput := K8sNodeInput{
		Name:          "master",
		Node:          1,
		DiskSize:      1,
		Flist:         "",
		CPU:           1,
		Memory:        1,
		PublicIP:      false,
		PublicIP6:     false,
		Planetary:     false,
		FlistChecksum: "checksum",
	}

	k8sWorkerInput1 := k8sNodeInput
	k8sWorkerInput1.Name = "worker1"

	k8sWorkerInput2 := k8sNodeInput
	k8sWorkerInput2.Name = "worker2"

	k8sInput := KubernetesArgs{
		Master:       k8sNodeInput,
		Workers:      []K8sNodeInput{k8sWorkerInput1, k8sWorkerInput2},
		Token:        "token",
		NetworkName:  "network",
		SolutionType: "solution",
		SSHKey:       "ssh",
	}

	t.Run("parsing input success", func(t *testing.T) {
		k8s, err := parseToK8sCluster(k8sInput)
		assert.NoError(t, err)
		assert.Equal(t, k8s.Master.Name, k8sInput.Master.Name)
		assert.Equal(t, len(k8s.Workers), len(k8sInput.Workers))
		assert.Equal(t, k8s.Token, k8sInput.Token)
		assert.Equal(t, k8s.NetworkName, k8sInput.NetworkName)
		assert.Equal(t, k8s.SolutionType, k8sInput.SolutionType)
		assert.Equal(t, k8s.SSHKey, k8sInput.SSHKey)
	})

	t.Run("parsing input failed: wrong node id type", func(t *testing.T) {
		k8sInput.Master.Node = ""
		_, err := parseToK8sCluster(k8sInput)
		assert.Error(t, err)
		k8sInput.Master.Node = 1
	})

	t.Run("parsing input failed: wrong worker node id type", func(t *testing.T) {
		k8sInput.Workers[0].Node = ""
		_, err := parseToK8sCluster(k8sInput)
		assert.Error(t, err)
		k8sInput.Workers[0].Node = 1
	})

	t.Run("parsing and update k8s success", func(t *testing.T) {
		k8s, err := parseToK8sCluster(k8sInput)
		assert.NoError(t, err)

		k8s.NodeDeploymentID = map[uint32]uint64{1: 1}
		ip, err := gridtypes.ParseIPNet("1.1.1.1/16")
		assert.NoError(t, err)
		k8s.NodesIPRange = map[uint32]gridtypes.IPNet{1: ip}

		state := parseToK8sState(k8s)
		assert.Equal(t, k8s.NodeDeploymentID[1], uint64(state.NodeDeploymentID["1"]))
		assert.Equal(t, k8s.NodesIPRange[1].String(), state.NodesIPRange["1"])

		err = updateK8sFromState(&k8s, state)
		assert.NoError(t, err)
	})

	t.Run("parsing and update k8s failed: wrong node id type (node deployment)", func(t *testing.T) {
		k8s, err := parseToK8sCluster(k8sInput)
		assert.NoError(t, err)

		state := parseToK8sState(k8s)
		state.NodeDeploymentID = map[string]int64{"": 1}
		err = updateK8sFromState(&k8s, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update k8s failed: wrong node id type (node ip ranges)", func(t *testing.T) {
		k8s, err := parseToK8sCluster(k8sInput)
		assert.NoError(t, err)

		state := parseToK8sState(k8s)
		state.NodesIPRange = map[string]string{"": "1.1.1.1/16"}
		err = updateK8sFromState(&k8s, state)
		assert.Error(t, err)
	})

	t.Run("parsing and update k8s failed: wrong node ip ranges type", func(t *testing.T) {
		k8s, err := parseToK8sCluster(k8sInput)
		assert.NoError(t, err)

		state := parseToK8sState(k8s)
		state.NodesIPRange = map[string]string{"": "1"}
		err = updateK8sFromState(&k8s, state)
		assert.Error(t, err)
	})
}
