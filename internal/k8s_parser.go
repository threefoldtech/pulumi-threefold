package provider

import (
	"fmt"
	"strconv"

	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

// K8sNodeInput struct of input data
type K8sNodeInput struct {
	Name          string `pulumi:"name"`
	Node          int32  `pulumi:"node"`
	DiskSize      int    `pulumi:"disk_size"`
	Flist         string `pulumi:"flist,optional"`
	CPU           int    `pulumi:"cpu"`
	Memory        int    `pulumi:"memory"`
	PublicIP      bool   `pulumi:"publicip,optional"`
	PublicIP6     bool   `pulumi:"publicip6,optional"`
	Planetary     bool   `pulumi:"planetary,optional"`
	FlistChecksum string `pulumi:"flist_checksum,optional"`
}

// K8sNodeComputed struct of computed data
type K8sNodeComputed struct {
	ComputedIP  string `pulumi:"computedip"`
	ComputedIP6 string `pulumi:"computedip6"`
	IP          string `pulumi:"ip"`
	YggIP       string `pulumi:"ygg_ip"`
	ConsoleURL  string `pulumi:"console_url"`
	SSHKey      string `pulumi:"ssh_key"`
	Token       string `pulumi:"token"`
	NetworkName string `pulumi:"network_name"`
}

func parseToKubernetesState(k8sCluster workloads.K8sCluster) KubernetesState {

	// parse NodesIpRange
	nodesIPRange := make(map[string]string)
	for k, v := range k8sCluster.NodesIPRange {
		nodesIPRange[fmt.Sprint(k)] = v.String()
	}

	// parse NodeDeploymentID
	nodeDeploymentID := make(map[string]int64)
	for k, v := range k8sCluster.NodeDeploymentID {
		nodeDeploymentID[fmt.Sprint(k)] = int64(v)
	}

	// parse master computed
	masterComputed := &K8sNodeComputed{
		ComputedIP:  k8sCluster.Master.ComputedIP,
		ComputedIP6: k8sCluster.Master.ComputedIP6,
		YggIP:       k8sCluster.Master.YggIP,
		IP:          k8sCluster.Master.IP,
		NetworkName: k8sCluster.Master.NetworkName,
		Token:       k8sCluster.Master.Token,
		SSHKey:      k8sCluster.Master.SSHKey,
		ConsoleURL:  k8sCluster.Master.ConsoleURL,
	}

	// parse master input
	masterInput := &K8sNodeInput{
		Name:          k8sCluster.Master.Name,
		Node:          int32(k8sCluster.Master.Node),
		DiskSize:      k8sCluster.Master.DiskSize,
		PublicIP:      k8sCluster.Master.PublicIP,
		PublicIP6:     k8sCluster.Master.PublicIP6,
		Planetary:     k8sCluster.Master.Planetary,
		Flist:         k8sCluster.Master.Flist,
		FlistChecksum: k8sCluster.Master.FlistChecksum,
		CPU:           k8sCluster.Master.CPU,
		Memory:        k8sCluster.Master.Memory,
	}

	// parse workers computed & input
	workersComputed := make(map[string]K8sNodeComputed)
	workersInput := []K8sNodeInput{}
	for _, w := range k8sCluster.Workers {
		newWorkerComputed := K8sNodeComputed{
			ComputedIP:  w.ComputedIP,
			ComputedIP6: w.ComputedIP6,
			YggIP:       w.YggIP,
			IP:          w.IP,
			NetworkName: w.NetworkName,
			Token:       w.Token,
			SSHKey:      w.SSHKey,
			ConsoleURL:  w.ConsoleURL,
		}
		workersComputed[w.Name] = newWorkerComputed

		newWorkerInput := K8sNodeInput{
			Name:          w.Name,
			Node:          int32(w.Node),
			DiskSize:      w.DiskSize,
			PublicIP:      w.PublicIP,
			PublicIP6:     w.PublicIP6,
			Planetary:     w.Planetary,
			Flist:         w.Flist,
			FlistChecksum: w.FlistChecksum,
			CPU:           w.CPU,
			Memory:        w.Memory,
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

func parseToK8sCluster(kubernetesArgs KubernetesArgs) workloads.K8sCluster {

	// parse master
	master := &workloads.K8sNode{
		Name:          kubernetesArgs.Master.Name,
		Node:          uint32(kubernetesArgs.Master.Node),
		DiskSize:      kubernetesArgs.Master.DiskSize,
		PublicIP:      kubernetesArgs.Master.PublicIP,
		PublicIP6:     kubernetesArgs.Master.PublicIP6,
		Planetary:     kubernetesArgs.Master.Planetary,
		Flist:         kubernetesArgs.Master.Flist,
		FlistChecksum: kubernetesArgs.Master.FlistChecksum,
		CPU:           kubernetesArgs.Master.CPU,
		Memory:        kubernetesArgs.Master.Memory,
	}

	// set default flist
	if master.Flist == "" {
		master.Flist = "https://hub.grid.tf/tf-official-apps/threefoldtech-k3s-latest.flist"
	}

	// parse workers
	workers := []workloads.K8sNode{}
	for _, w := range kubernetesArgs.Workers {
		newWorker := workloads.K8sNode{
			Name:          w.Name,
			Node:          uint32(w.Node),
			DiskSize:      w.DiskSize,
			PublicIP:      w.PublicIP,
			PublicIP6:     w.PublicIP6,
			Planetary:     w.Planetary,
			Flist:         w.Flist,
			FlistChecksum: w.FlistChecksum,
			CPU:           w.CPU,
			Memory:        w.Memory,
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
	}
}

func addComputedFieldsFromState(k8sCluster *workloads.K8sCluster, state KubernetesState) error {

	// parse NodesIPRange
	nodesIPRange := make(map[uint32]gridtypes.IPNet)
	for k, v := range state.NodesIPRange {
		ipRange, err := gridtypes.ParseIPNet(v)
		if err != nil {
			return err
		}

		kInt, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		nodesIPRange[uint32(kInt)] = ipRange
	}

	// update NodesIPRange
	k8sCluster.NodesIPRange = nodesIPRange

	// parse NodeDeploymentID
	nodeDeploymentID := make(map[uint32]uint64)
	for k, v := range state.NodeDeploymentID {
		kInt, err := strconv.Atoi(k)
		if err != nil {
			return err
		}

		nodeDeploymentID[uint32(kInt)] = uint64(v)
	}

	// update NodeDeploymentID
	k8sCluster.NodeDeploymentID = nodeDeploymentID

	// update master computed
	k8sCluster.Master.ComputedIP = state.MasterComputed.ComputedIP
	k8sCluster.Master.ComputedIP6 = state.MasterComputed.ComputedIP6
	k8sCluster.Master.IP = state.MasterComputed.IP
	k8sCluster.Master.YggIP = state.MasterComputed.YggIP
	k8sCluster.Master.ConsoleURL = state.MasterComputed.ConsoleURL
	k8sCluster.Master.SSHKey = state.MasterComputed.SSHKey
	k8sCluster.Master.Token = state.MasterComputed.Token
	k8sCluster.Master.NetworkName = state.MasterComputed.NetworkName

	// update workers computed
	for _, v := range k8sCluster.Workers {
		// update every worker in k8sCluster if it has computed data in the state

		if workerComputed, ok := state.WorkersComputed[v.Name]; ok {
			v.ComputedIP = workerComputed.ComputedIP
			v.ComputedIP6 = workerComputed.ComputedIP6
			v.IP = workerComputed.IP
			v.YggIP = workerComputed.YggIP
			v.ConsoleURL = workerComputed.ConsoleURL
			v.SSHKey = workerComputed.SSHKey
			v.Token = workerComputed.Token
			v.NetworkName = workerComputed.NetworkName
		}
	}

	return nil
}
