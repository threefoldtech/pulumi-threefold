package provider

import (
	"context"
	"os"
	"strings"

	"github.com/pkg/errors"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi-go-provider/middleware/schema"
	goGen "github.com/pulumi/pulumi/pkg/v3/codegen/go"
	nodejsGen "github.com/pulumi/pulumi/pkg/v3/codegen/nodejs"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/threefoldtech/tfgrid-sdk-go/grid-client/deployer"
)

// Provider data and resources
func Provider() p.Provider {
	return infer.Provider(infer.Options{
		Metadata: schema.Metadata{
			Description:       "The Pulumi Resource Provider for the Threefold Grid.",
			DisplayName:       "Threefold Grid",
			Keywords:          []string{"pulumi", "grid", "threefold", "category/infrastructure", "kind/native"},
			Homepage:          "https://www.pulumi.com",
			Repository:        "https://github.com/threefoldtech/pulumi-threefold",
			Publisher:         "Threefold",
			LogoURL:           "https://www.threefold.io/images/black_threefold.png",
			License:           "Apache-2.0",
			PluginDownloadURL: "github://api.github.com/threefoldtech/pulumi-threefold",
			LanguageMap: map[string]any{
				"go": goGen.GoPackageInfo{
					GenerateExtraInputTypes:        true,
					GenerateResourceContainerTypes: true,
					ImportBasePath:                 "github.com/threefoldtech/pulumi-threefold/sdk/go/threefold",
				},
				"nodejs": nodejsGen.NodePackageInfo{
					PackageName: "@threefold/pulumi",
				},
			},
		},
		Resources: []infer.InferredResource{
			infer.Resource[*Scheduler](),
			infer.Resource[*Network](),
			infer.Resource[*Deployment](),
			infer.Resource[*Kubernetes](),
			infer.Resource[*GatewayName](),
			infer.Resource[*GatewayFQDN](),
		},
		Config: infer.Config[*Config](),
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		},
	})
}

// RunProvider runs the pulumi provider and adds its resources
func RunProvider(providerName, Version string) error {
	return p.RunProvider(providerName, Version, Provider())
}

// Config struct holds the configuration fields for the provider
type Config struct {
	Mnemonic     string   `pulumi:"mnemonic,optional"  provider:"secret"`
	Network      string   `pulumi:"network,optional"`
	KeyType      string   `pulumi:"key_type,optional"`
	SubstrateURL string   `pulumi:"substrate_url,optional"`
	RelayURLs    []string `pulumi:"relay_url,optional"`
	RmbTimeout   string   `pulumi:"rmb_timeout,optional"`

	TFPluginClient deployer.TFPluginClient
}

var _ = (infer.Annotated)((*Config)(nil))

// Annotate sets description and default values for configs
func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.Mnemonic, "The mnemonic of the user. It is very secret.")
	a.Describe(&c.Network, "The network to deploy on.")
	a.Describe(&c.KeyType, "The key type registered on substrate (ed25519 or sr25519).")
	a.Describe(&c.SubstrateURL, "The substrate url, example: wss://tfchain.dev.grid.tf/ws")
	a.Describe(&c.RelayURLs, "The relay urls, example: wss://relay.dev.grid.tf")
	a.Describe(&c.RmbTimeout, "The timeout duration in seconds for rmb calls")
	a.SetDefault(&c.Mnemonic, os.Getenv("MNEMONIC"), "")
	a.SetDefault(&c.Network, os.Getenv("NETWORK"), "")
	a.SetDefault(&c.KeyType, "sr25519", "")
}

var _ = (infer.CustomConfigure)((*Config)(nil))

// Configure checks configuration for the provider
func (c *Config) Configure(ctx context.Context) error {
	if len(strings.TrimSpace(c.Mnemonic)) == 0 {
		return errors.New("mnemonic is required")
	}

	if len(strings.TrimSpace(c.Network)) == 0 {
		c.Network = "dev"
	}

	opts := []deployer.PluginOpt{
		deployer.WithKeyType(c.KeyType),
		deployer.WithNetwork(c.Network),
	}

	if c.SubstrateURL != "" {
		opts = append(opts, deployer.WithSubstrateURL(c.SubstrateURL))
	}

	if len(c.RelayURLs) > 0 {
		opts = append(opts, deployer.WithRelayURL(c.RelayURLs...))
	}

	tfPluginClient, err := deployer.NewTFPluginClient(
		c.Mnemonic,
		opts...,
	)
	if err != nil {
		return errors.Wrap(err, "error creating threefold plugin client")
	}

	c.TFPluginClient = tfPluginClient

	p.GetLogger(ctx).Info("threefold grid provider setup")

	return nil
}
