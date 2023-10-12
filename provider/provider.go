package provider

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
)

// RunProvider runs the pulumi provider and adds its resources
func RunProvider(providerName, Version string) error {
	return p.RunProvider(providerName, Version,
		infer.Provider(infer.Options{
			Resources: []infer.InferredResource{
				infer.Resource[*Scheduler, SchedulerArgs, SchedulerState](),
				infer.Resource[*Network, NetworkArgs, NetworkState](),
				infer.Resource[*Deployment, DeploymentArgs, DeploymentState](),
				infer.Resource[*Kubernetes, KubernetesArgs, KubernetesState](),
				infer.Resource[*GatewayName, GatewayNameArgs, GatewayNameState](),
				infer.Resource[*GatewayFQDN, GatewayFQDNArgs, GatewayFQDNState](),
			},
			Config: infer.Config[*Config](),
		}))
}

// Config struct holds the configuration fields for the provider
type Config struct {
	Mnemonic     string `pulumi:"mnemonic,optional"  provider:"secret"`
	Network      string `pulumi:"network,optional"`
	KeyType      string `pulumi:"key_type,optional"`
	SubstrateURL string `pulumi:"substrate_url,optional"`
	RelayURL     string `pulumi:"relay_url,optional"`
	RmbTimeout   string `pulumi:"rmb_timeout,optional"`

	TFPluginClient deployer.TFPluginClient
}

var _ = (infer.Annotated)((*Config)(nil))

// Annotate sets description and default values for configs
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.Mnemonic, "The mnemonic of the user. It is very secret.")
	a.Describe(&c.Network, "The network to deploy on.")
	a.Describe(&c.KeyType, "The key type registered on substrate (ed25519 or sr25519).")
	a.Describe(&c.SubstrateURL, "The substrate url, example: wss://tfchain.dev.grid.tf/ws")
	a.Describe(&c.RelayURL, "The relay url, example: wss://relay.dev.grid.tf")
	a.Describe(&c.RmbTimeout, "The timeout duration in seconds for rmb calls")
	a.SetDefault(&c.Mnemonic, os.Getenv("MNEMONIC"), "")
	a.SetDefault(&c.Network, os.Getenv("NETWORK"), "")
	a.SetDefault(&c.KeyType, "sr25519", "")
}

var _ = (infer.CustomConfigure)((*Config)(nil))

// Configure checks configuration for the provider
func (c *Config) Configure(ctx p.Context) error {
	if len(strings.TrimSpace(c.Mnemonic)) == 0 {
		return errors.New("mnemonic is required")
	}

	if len(strings.TrimSpace(c.Network)) == 0 {
		c.Network = "dev"
	}

	tfPluginClient, err := deployer.NewTFPluginClient(
		c.Mnemonic, c.KeyType, c.Network, c.SubstrateURL, c.RelayURL, "", 0, false,
	)
	if err != nil {
		return errors.Wrap(err, "error creating threefold plugin client")
	}

	c.TFPluginClient = tfPluginClient

	ctx.Log(diag.Info, "threefold grid provider setup")

	return nil
}
