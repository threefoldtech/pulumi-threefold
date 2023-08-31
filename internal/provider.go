package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
)

func RunProvider(providerName, Version string) error {
	return p.RunProvider(providerName, Version,
		infer.Provider(infer.Options{
			Resources: []infer.InferredResource{},
		}))
}
