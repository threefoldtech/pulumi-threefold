package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

// GatewayName controlling struct
type GatewayName struct{}

// GatewayNameArgs is defining what arguments it accepts
type GatewayNameArgs struct {
	Name           string   `pulumi:"name"`
	NodeID         int32    `pulumi:"node_id"`
	Backends       []string `pulumi:"backends"`
	TLSPassthrough bool     `pulumi:"tls_passthrough,optional"`
	Network        string   `pulumi:"network,optional"`
	Description    string   `pulumi:"description,optional"`
	SolutionType   string   `pulumi:"solution_type,optional"`
}

// GatewayNameState is describing the fields that exist on the created resource.
type GatewayNameState struct {
	GatewayNameArgs

	NodeDeploymentID map[int32]int64 `pulumi:"node_deployment_id"`
	FQDN             string          `pulumi:"fqdn"`
	NameContractID   int64           `pulumi:"name_contract_id"`
	ContractID       int64           `pulumi:"contract_id"`
}

// Create creates GatewayName and deploy it
func (*GatewayName) Create(ctx p.Context, name string, input GatewayNameArgs, preview bool) (string, GatewayNameState, error) {
	state := GatewayNameState{GatewayNameArgs: input}
	if preview {
		return name, state, nil
	}

	// k8sCluster := parseToK8sCluster(input)

	// config := infer.GetConfig[Config](ctx)

	// if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
	// 	return name, state, err
	// }

	// if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
	// 	return name, state, err
	// }

	// state = parseToKubernetesState(k8sCluster)

	return name, state, nil
}

// Update updates the GatewayName resource
func (*GatewayName) Update(ctx p.Context, id string, oldState GatewayNameState, input GatewayNameArgs, preview bool) (GatewayNameState, error) {
	state := GatewayNameState{GatewayNameArgs: input}
	if preview {
		return state, nil
	}

	// k8sCluster := parseToK8sCluster(input)
	// if err := addComputedFieldsFromState(&k8sCluster, oldState); err != nil {
	// 	return state, err
	// }

	// config := infer.GetConfig[Config](ctx)

	// if err := config.TFPluginClient.K8sDeployer.Deploy(ctx, &k8sCluster); err != nil {
	// 	return state, err
	// }

	// if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
	// 	return state, err
	// }

	// state = parseToKubernetesState(k8sCluster)

	return state, nil
}

// Read gets the state of the GatewayName resource
// func (*GatewayName) Read(ctx p.Context, id string, oldState GatewayNameState) (string, GatewayNameState, error) {
	// k8sCluster := parseToK8sCluster(oldState.GatewayNameArgs)
	// if err := addComputedFieldsFromState(&k8sCluster, oldState); err != nil {
	// 	return id, oldState, err
	// }

	// config := infer.GetConfig[Config](ctx)

	// if err := config.TFPluginClient.K8sDeployer.Validate(ctx, &k8sCluster); err != nil {
	// 	return id, oldState, err
	// }

	// if err := k8sCluster.InvalidateBrokenAttributes(config.TFPluginClient.SubstrateConn); err != nil {
	// 	return id, oldState, err
	// }

	// if err := config.TFPluginClient.K8sDeployer.UpdateFromRemote(ctx, &k8sCluster); err != nil {
	// 	return id, oldState, err
	// }

	// state := parseToKubernetesState(k8sCluster)

	// return id, state, nil
// }

// Delete deletes the GatewayName resource
// func (*GatewayName) Delete(ctx p.Context, id string, oldState GatewayNameState) error {
// 	k8sCluster := parseToK8sCluster(oldState.GatewayNameArgs)
// 	if err := addComputedFieldsFromState(&k8sCluster, oldState); err != nil {
// 		return err
// 	}

// 	config := infer.GetConfig[Config](ctx)

// 	return config.TFPluginClient.K8sDeployer.Cancel(ctx, &k8sCluster)
// }
