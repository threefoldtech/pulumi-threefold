package provider

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

// K8sNodeInput struct of input data
type K8sNodeInput struct {
	Name           string      `pulumi:"name"`
	NetworkName    string      `pulumi:"network_name"`
	Node           interface{} `pulumi:"node"`
	DiskSize       int         `pulumi:"disk_size"`
	Flist          string      `pulumi:"flist,optional"`
	FlistChecksum  string      `pulumi:"flist_checksum,optional"`
	CPU            int         `pulumi:"cpu"`
	Memory         int         `pulumi:"memory"`
	PublicIP       bool        `pulumi:"public_ip,optional"`
	PublicIP6      bool        `pulumi:"public_ip6,optional"`
	Planetary      bool        `pulumi:"planetary,optional"`
	Mycelium       bool        `pulumi:"mycelium,optional"`
	MyceliumIPSeed string      `pulumi:"mycelium_ip_seed,optional"`
}

// K8sNodeComputed struct of computed data
type K8sNodeComputed struct {
	MyceliumIPSeed string `pulumi:"mycelium_ip_seed"`
	ComputedIP     string `pulumi:"computed_ip"`
	ComputedIP6    string `pulumi:"computed_ip6"`
	IP             string `pulumi:"ip"`
	PlanetaryIP    string `pulumi:"planetary_ip"`
	MyceliumIP     string `pulumi:"mycelium_ip"`
	ConsoleURL     string `pulumi:"console_url"`
}

func parseToK8sState(k8sCluster workloads.K8sCluster) KubernetesState {
	// parse NodesIpRange
	nodesIPRange := make(map[string]string)
	for nodeID, ipRange := range k8sCluster.NodesIPRange {
		nodesIPRange[fmt.Sprint(nodeID)] = ipRange.String()
	}

	// parse NodeDeploymentID
	nodeDeploymentID := make(map[string]int64)
	for nodeID, deploymentID := range k8sCluster.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(nodeID)] = int64(deploymentID)
	}

	// parse master computed
	masterComputed := K8sNodeComputed{
		MyceliumIPSeed: hex.EncodeToString(k8sCluster.Master.MyceliumIPSeed),
		ComputedIP:     k8sCluster.Master.ComputedIP,
		ComputedIP6:    k8sCluster.Master.ComputedIP6,
		PlanetaryIP:    k8sCluster.Master.PlanetaryIP,
		MyceliumIP:     k8sCluster.Master.MyceliumIP,
		IP:             k8sCluster.Master.IP,
		ConsoleURL:     k8sCluster.Master.ConsoleURL,
	}

	// parse master input
	masterInput := K8sNodeInput{
		Name:           k8sCluster.Master.Name,
		NetworkName:    k8sCluster.Master.NetworkName,
		Node:           int(k8sCluster.Master.NodeID),
		DiskSize:       int(k8sCluster.Master.DiskSizeGB),
		PublicIP:       k8sCluster.Master.PublicIP,
		PublicIP6:      k8sCluster.Master.PublicIP6,
		Planetary:      k8sCluster.Master.Planetary,
		MyceliumIPSeed: hex.EncodeToString(k8sCluster.Master.MyceliumIPSeed),
		Flist:          k8sCluster.Master.Flist,
		FlistChecksum:  k8sCluster.Master.FlistChecksum,
		CPU:            int(k8sCluster.Master.CPU),
		Memory:         int(k8sCluster.Master.MemoryMB),
	}

	// parse workers computed & input
	workersComputed := make(map[string]K8sNodeComputed)
	workersInput := []K8sNodeInput{}
	for _, w := range k8sCluster.Workers {
		newWorkerComputed := K8sNodeComputed{
			MyceliumIPSeed: hex.EncodeToString(w.MyceliumIPSeed),
			ComputedIP:     w.ComputedIP,
			ComputedIP6:    w.ComputedIP6,
			PlanetaryIP:    w.PlanetaryIP,
			MyceliumIP:     w.MyceliumIP,
			IP:             w.IP,
			ConsoleURL:     w.ConsoleURL,
		}
		workersComputed[w.Name] = newWorkerComputed

		newWorkerInput := K8sNodeInput{
			Name:           w.Name,
			NetworkName:    w.NetworkName,
			Node:           int(w.NodeID),
			DiskSize:       int(w.DiskSizeGB),
			PublicIP:       w.PublicIP,
			PublicIP6:      w.PublicIP6,
			Planetary:      w.Planetary,
			MyceliumIPSeed: hex.EncodeToString(w.MyceliumIPSeed),
			Flist:          w.Flist,
			FlistChecksum:  w.FlistChecksum,
			CPU:            int(w.CPU),
			Memory:         int(w.MemoryMB),
		}
		workersInput = append(workersInput, newWorkerInput)
	}

	return KubernetesState{
		KubernetesArgs: KubernetesArgs{
			Master:       masterInput,
			Workers:      workersInput,
			Token:        k8sCluster.Token,
			NetworkName:  k8sCluster.NetworkName,
			SolutionType: k8sCluster.SolutionType,
			SSHKey:       k8sCluster.SSHKey,
		},
		MasterComputed:   masterComputed,
		WorkersComputed:  workersComputed,
		NodesIPRange:     nodesIPRange,
		NodeDeploymentID: nodeDeploymentID,
	}
}

func parseToK8sCluster(kubernetesArgs KubernetesArgs) (workloads.K8sCluster, error) {
	// for tests
	sshKey := os.Getenv("SSH_KEY")
	if sshKey != "" {
		kubernetesArgs.SSHKey = sshKey
	}

	nodeID, err := strconv.Atoi(fmt.Sprint(kubernetesArgs.Master.Node))
	if err != nil {
		return workloads.K8sCluster{}, err
	}

	myceliumIPSeed, err := hex.DecodeString(kubernetesArgs.Master.MyceliumIPSeed)
	if err != nil {
		return workloads.K8sCluster{}, err
	}

	if kubernetesArgs.Master.Mycelium && len(kubernetesArgs.Master.MyceliumIPSeed) == 0 {
		myceliumIPSeed, err = workloads.RandomMyceliumIPSeed()
		if err != nil {
			return workloads.K8sCluster{}, err
		}
	}

	if len(kubernetesArgs.Master.NetworkName) == 0 {
		kubernetesArgs.Master.NetworkName = kubernetesArgs.NetworkName
	}

	// parse master
	master := &workloads.K8sNode{
		VM: &workloads.VM{
			Name:           kubernetesArgs.Master.Name,
			NetworkName:    kubernetesArgs.Master.NetworkName,
			NodeID:         uint32(nodeID),
			PublicIP:       kubernetesArgs.Master.PublicIP,
			PublicIP6:      kubernetesArgs.Master.PublicIP6,
			Planetary:      kubernetesArgs.Master.Planetary,
			MyceliumIPSeed: myceliumIPSeed,
			Flist:          kubernetesArgs.Master.Flist,
			FlistChecksum:  kubernetesArgs.Master.FlistChecksum,
			CPU:            uint8(kubernetesArgs.Master.CPU),
			MemoryMB:       uint64(kubernetesArgs.Master.Memory),
		},
		DiskSizeGB: uint64(kubernetesArgs.Master.DiskSize),
	}

	// set default flist
	if master.Flist == "" {
		master.Flist = "https://hub.grid.tf/tf-official-apps/threefoldtech-k3s-latest.flist"
	}

	// parse workers
	workers := []workloads.K8sNode{}
	for _, w := range kubernetesArgs.Workers {
		nodeID, err := strconv.Atoi(fmt.Sprint(w.Node))
		if err != nil {
			return workloads.K8sCluster{}, err
		}

		myceliumIPSeed, err := hex.DecodeString(w.MyceliumIPSeed)
		if err != nil {
			return workloads.K8sCluster{}, err
		}

		if w.Mycelium && len(w.MyceliumIPSeed) == 0 {
			myceliumIPSeed, err = workloads.RandomMyceliumIPSeed()
			if err != nil {
				return workloads.K8sCluster{}, err
			}
		}

		if len(kubernetesArgs.Master.NetworkName) == 0 {
			w.NetworkName = kubernetesArgs.NetworkName
		}

		newWorker := workloads.K8sNode{
			VM: &workloads.VM{
				Name:           w.Name,
				NetworkName:    w.NetworkName,
				NodeID:         uint32(nodeID),
				PublicIP:       w.PublicIP,
				PublicIP6:      w.PublicIP6,
				Planetary:      w.Planetary,
				MyceliumIPSeed: myceliumIPSeed,
				Flist:          w.Flist,
				FlistChecksum:  w.FlistChecksum,
				CPU:            uint8(w.CPU),
				MemoryMB:       uint64(w.Memory),
			},
			DiskSizeGB: uint64(w.DiskSize),
		}

		// set default flist
		if newWorker.Flist == "" {
			newWorker.Flist = "https://hub.grid.tf/tf-official-apps/threefoldtech-k3s-latest.flist"
		}

		workers = append(workers, newWorker)
	}

	return workloads.K8sCluster{
		Master:       master,
		Workers:      workers,
		Token:        kubernetesArgs.Token,
		NetworkName:  kubernetesArgs.NetworkName,
		SolutionType: kubernetesArgs.SolutionType,
		SSHKey:       kubernetesArgs.SSHKey,
		NodesIPRange: make(map[uint32]gridtypes.IPNet),
	}, nil
}

func updateK8sFromState(k8sCluster *workloads.K8sCluster, state KubernetesState) error {
	// parse NodesIPRange
	nodesIPRange := make(map[uint32]gridtypes.IPNet)
	for nodeID, ipRange := range state.NodesIPRange {
		ipRange, err := gridtypes.ParseIPNet(ipRange)
		if err != nil {
			return err
		}

		node, err := strconv.Atoi(nodeID)
		if err != nil {
			return err
		}

		nodesIPRange[uint32(node)] = ipRange
	}

	// update NodesIPRange
	k8sCluster.NodesIPRange = nodesIPRange

	// parse NodeDeploymentID
	nodeDeploymentID := make(map[uint32]uint64)
	for nodeID, deploymentID := range state.NodeDeploymentID {
		node, err := strconv.Atoi(nodeID)
		if err != nil {
			return err
		}

		nodeDeploymentID[uint32(node)] = uint64(deploymentID)
	}

	// update NodeDeploymentID
	k8sCluster.NodeDeploymentID = nodeDeploymentID

	myceliumIPSeed, err := hex.DecodeString(state.MasterComputed.MyceliumIPSeed)
	if err != nil {
		return err
	}

	// update master computed
	k8sCluster.Master.ComputedIP = state.MasterComputed.ComputedIP
	k8sCluster.Master.ComputedIP6 = state.MasterComputed.ComputedIP6
	k8sCluster.Master.IP = state.MasterComputed.IP
	k8sCluster.Master.PlanetaryIP = state.MasterComputed.PlanetaryIP
	k8sCluster.Master.MyceliumIP = state.MasterComputed.MyceliumIP
	k8sCluster.Master.ConsoleURL = state.MasterComputed.ConsoleURL
	k8sCluster.Master.MyceliumIPSeed = myceliumIPSeed

	// update workers computed
	for i, worker := range k8sCluster.Workers {
		// update every worker in k8sCluster if it has computed data in the state
		if workerComputed, ok := state.WorkersComputed[worker.Name]; ok {
			myceliumIPSeed, err := hex.DecodeString(workerComputed.MyceliumIPSeed)
			if err != nil {
				return err
			}

			k8sCluster.Workers[i].ComputedIP = workerComputed.ComputedIP
			k8sCluster.Workers[i].ComputedIP6 = workerComputed.ComputedIP6
			k8sCluster.Workers[i].IP = workerComputed.IP
			k8sCluster.Workers[i].PlanetaryIP = workerComputed.PlanetaryIP
			k8sCluster.Workers[i].MyceliumIP = workerComputed.MyceliumIP
			k8sCluster.Workers[i].ConsoleURL = workerComputed.ConsoleURL
			k8sCluster.Workers[i].MyceliumIPSeed = myceliumIPSeed
		}
	}

	return nil
}
