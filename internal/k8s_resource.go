package provider

import (
	"fmt"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

// Kubernetes struct data
type K8sNode struct {
	Name          string `pulumi:"name"`
	Node          int32  `pulumi:"node"`
	DiskSize      int    `pulumi:"disk_size"`
	PublicIP      bool   `pulumi:"publicip,optional"`
	PublicIP6     bool   `pulumi:"publicip6,optional"`
	Planetary     bool   `pulumi:"planetary,optional"`
	Flist         string `pulumi:"flist,optional"`
	FlistChecksum string `pulumi:"flist_checksum,optional"`
	ComputedIP    string `pulumi:"computedip,optional"`
	ComputedIP6   string `pulumi:"computedip6,optional"`
	YggIP         string `pulumi:"ygg_ip,optional"`
	IP            string `pulumi:"ip,optional"`
	CPU           int    `pulumi:"cpu"`
	Memory        int    `pulumi:"memory"`
	NetworkName   string `pulumi:"network_name,optional"`
	Token         string `pulumi:"token,optional"`
	SSHKey        string `pulumi:"ssh_key,optional"`
	ConsoleURL    string `pulumi:"console_url,optional"`
}

// Kubernetes controlling struct
type Kubernetes struct{}

// KubernetesArgs is defining what arguments it accepts
type KubernetesArgs struct {
	Master       *K8sNode  `pulumi:"master"`
	Workers      []K8sNode `pulumi:"workers"`
	Token        string    `pulumi:"token"`
	NetworkName  string    `pulumi:"network_name"`
	SolutionType string    `pulumi:"solution_type,optional"`
	SSHKey       string    `pulumi:"ssh_key,optional"`
}

type KubernetesState struct {
	KubernetesArgs

	NodesIPRange     map[string]string `pulumi:"nodes_ip_range"`
	NodeDeploymentID map[string]int64  `pulumi:"node_deployment_id"`
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

	// parse master
	master := &K8sNode{
		Name:          k8sCluster.Master.Name,
		Node:          int32(k8sCluster.Master.Node),
		DiskSize:      k8sCluster.Master.DiskSize,
		PublicIP:      k8sCluster.Master.PublicIP,
		PublicIP6:     k8sCluster.Master.PublicIP6,
		Planetary:     k8sCluster.Master.Planetary,
		Flist:         k8sCluster.Master.Flist,
		FlistChecksum: k8sCluster.Master.FlistChecksum,
		ComputedIP:    k8sCluster.Master.ComputedIP,
		ComputedIP6:   k8sCluster.Master.ComputedIP6,
		YggIP:         k8sCluster.Master.YggIP,
		IP:            k8sCluster.Master.IP,
		CPU:           k8sCluster.Master.CPU,
		Memory:        k8sCluster.Master.Memory,
		NetworkName:   k8sCluster.Master.NetworkName,
		Token:         k8sCluster.Master.Token,
		SSHKey:        k8sCluster.Master.SSHKey,
		ConsoleURL:    k8sCluster.Master.ConsoleURL,
	}

	// parse workers
	workers := []K8sNode{}
	for _, w := range k8sCluster.Workers {
		newWorker := K8sNode{
			Name:          w.Name,
			Node:          int32(w.Node),
			DiskSize:      w.DiskSize,
			PublicIP:      w.PublicIP,
			PublicIP6:     w.PublicIP6,
			Planetary:     w.Planetary,
			Flist:         w.Flist,
			FlistChecksum: w.FlistChecksum,
			ComputedIP:    w.ComputedIP,
			ComputedIP6:   w.ComputedIP6,
			YggIP:         w.YggIP,
			IP:            w.IP,
			CPU:           w.CPU,
			Memory:        w.Memory,
			NetworkName:   w.NetworkName,
			Token:         w.Token,
			SSHKey:        w.SSHKey,
			ConsoleURL:    w.ConsoleURL,
		}

		workers = append(workers, newWorker)
	}

	return KubernetesState{
		KubernetesArgs: KubernetesArgs{
			Master:       master,
			Workers:      workers,
			Token:        k8sCluster.Token,
			NetworkName:  k8sCluster.NetworkName,
			SolutionType: k8sCluster.SolutionType,
			SSHKey:       k8sCluster.SSHKey,
		},
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
		ComputedIP:    kubernetesArgs.Master.ComputedIP,
		ComputedIP6:   kubernetesArgs.Master.ComputedIP6,
		YggIP:         kubernetesArgs.Master.YggIP,
		IP:            kubernetesArgs.Master.IP,
		CPU:           kubernetesArgs.Master.CPU,
		Memory:        kubernetesArgs.Master.Memory,
		NetworkName:   kubernetesArgs.Master.NetworkName,
		Token:         kubernetesArgs.Master.Token,
		SSHKey:        kubernetesArgs.Master.SSHKey,
		ConsoleURL:    kubernetesArgs.Master.ConsoleURL,
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
			ComputedIP:    w.ComputedIP,
			ComputedIP6:   w.ComputedIP6,
			YggIP:         w.YggIP,
			IP:            w.IP,
			CPU:           w.CPU,
			Memory:        w.Memory,
			NetworkName:   w.NetworkName,
			Token:         w.Token,
			SSHKey:        w.SSHKey,
			ConsoleURL:    w.ConsoleURL,
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

// Create creates Kubernetes cluster and deploy it
func (*Kubernetes) Create(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster := parseToK8sCluster(input)

	config := infer.GetConfig[Config](ctx)
	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}

// Update updates the arguments of the Kubernetes resource
func (*Kubernetes) Update(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster := parseToK8sCluster(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}

// Read get the state of the Kubernetes resource
func (*Kubernetes) Read(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster := parseToK8sCluster(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Validate(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	if err := k8sCluster.InvalidateBrokenAttributes(config.TFPluginClient.SubstrateConn); err != nil {
		return name, state, err
	}

	if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}

// Delete deletes the Kubernetes resource
func (*Kubernetes) Delete(ctx p.Context, name string, input KubernetesArgs, preview bool) (string, KubernetesState, error) {
	state := KubernetesState{KubernetesArgs: input}
	if preview {
		return name, state, nil
	}

	k8sCluster := parseToK8sCluster(input)

	config := infer.GetConfig[Config](ctx)

	if err := config.TFPluginClient.K8sDeployer.Cancel(ctx, &k8sCluster); err != nil {
		return name, state, err
	}

	state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}
