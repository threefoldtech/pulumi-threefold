package provider

import (
	"github.com/pkg/errors"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/diag"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
)

func RunProvider(providerName, Version string) error {
	return p.RunProvider(providerName, Version,
		infer.Provider(infer.Options{
			Resources: []infer.InferredResource{
				infer.Resource[*Network, NetworkArgs, NetworkState](),
				infer.Resource[*Deployment, DeploymentArgs, DeploymentState](),
			},
			Config: infer.Config[*Config](),
		}))
}

type Config struct {
	Mnemonic      string `pulumi:"mnemonic"  provider:"secret"`
	Network       string `pulumi:"network,optional"`
	KeyType       string `pulumi:"key_type,optional"`
	Substrate_URL string `pulumi:"substrate_url,optional"`
	Relay_URL     string `pulumi:"relay_url,optional"`
	Rmb_Timeout   string `pulumi:"rmb_timeout,optional"`

	TFPluginClient deployer.TFPluginClient
}

var _ = (infer.Annotated)((*Config)(nil))

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.Mnemonic, "The mnemonic of the user. It is very secret.")
	a.Describe(&c.Network, "The network to deploy on.")
	a.Describe(&c.KeyType, "The key type registered on substrate (ed25519 or sr25519).")
	a.Describe(&c.Substrate_URL, "The substrate url, example: wss://tfchain.dev.grid.tf/ws")
	a.Describe(&c.Relay_URL, "The relay url, example: wss://relay.dev.grid.tf")
	a.Describe(&c.Rmb_Timeout, "The timeout duration in seconds for rmb calls")
	a.SetDefault(&c.Network, "", "dev")
	a.SetDefault(&c.KeyType, "", "sr25519")
}

var _ = (infer.CustomConfigure)((*Config)(nil))

func (c *Config) Configure(ctx p.Context) error {
	tfPluginClient, err := deployer.NewTFPluginClient(c.Mnemonic, c.KeyType, c.Network, c.Substrate_URL, c.Relay_URL, "", 0, false)
	if err != nil {
		return errors.Wrap(err, "error creating threefold plugin client")
	}

	c.TFPluginClient = tfPluginClient

	ctx.Log(diag.Info, "grid provider setup")

	return nil
}
