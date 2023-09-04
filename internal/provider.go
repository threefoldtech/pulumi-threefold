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
			},
			Config: infer.Config[*Config](),
		}))
}

type Config struct {
	Mnemonic string `pulumi:"mnemonic"  provider:"secret"`
	Network  string `pulumi:"network,optional"`
	KeyType  string `pulumi:"type,optional"`

	TFPluginClient deployer.TFPluginClient
}

// var _ = (infer.Annotated)((*Config)(nil))

// func (c *Config) Annotate(a infer.Annotator) {
// 	a.Describe(&c.Mnemonic, "The mnemonic of the user. It is very secret.")
// 	a.Describe(&c.Network, "The network to deploy on.")
// 	a.Describe(&c.keyType, "The network to deploy on.")
// 	a.SetDefault(&c.Network, "", "dev")
// 	a.SetDefault(&c.keyType, "", "sr25519")
// }

var _ = (infer.CustomConfigure)((*Config)(nil))

func (c *Config) Configure(ctx p.Context) error {
	tfPluginClient, err := deployer.NewTFPluginClient(c.Mnemonic, c.KeyType, c.Network, "", "", "", 0, false)
	if err != nil {
		return errors.Wrap(err, "error creating threefold plugin client")
	}

	c.TFPluginClient = tfPluginClient

	ctx.Log(diag.Info, "grid provider setup")

	return nil
}
