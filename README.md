# Grid Pulumi

[![Integration tests](https://github.com/threefoldtech/grid_pulumi/actions/workflows/integration.yaml/badge.svg?branch=development)](https://github.com/threefoldtech/grid_pulumi/actions/workflows/integration.yaml) [![Lint](https://github.com/threefoldtech/grid_pulumi/actions/workflows/lint.yaml/badge.svg?branch=development)](https://github.com/threefoldtech/grid_pulumi/actions/workflows/lint.yaml) [![Dependabot](https://badgen.net/badge/Dependabot/enabled/green?icon=dependabot)](https://dependabot.com/) <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-54%25-brightgreen.svg?longCache=true&style=flat)</a>

A Pulumi provider that enables infrastructure-as-code deployment of compute, network, and storage resources. It allows developers to provision virtual machines, networks, Kubernetes clusters, and storage using familiar Pulumi workflows in TypeScript, Python, Go, and .NET. The provider is available at the [Pulumi Registry](https://www.pulumi.com/registry/packages/threefold/).

## What this is

This provider maps grid resources into native Pulumi constructs, enabling declarative management of infrastructure through Pulumi's standard tooling. It supports resource lifecycle operations including create, update, and destroy, and integrates with the grid's identity and networking model.

## What this repository contains

- **Go-based Pulumi provider** implementing the resource bridge to grid APIs
- **Multi-language SDK generation** for TypeScript, Python, Go, and .NET
- **Example projects** demonstrating VM, network, Kubernetes, and storage deployments
- **Integration tests** validating provider functionality against live grid environments

## Role in the stack

The provider acts as the infrastructure-as-code entry point for the broader stack. It sits above the grid's deployment layer and allows users to define workloads declaratively rather than imperatively. It integrates with the grid's identity system (mnemonic-based keys) and network abstractions.

## Relation to ThreeFold

This technology is used within the ThreeFold ecosystem and was first deployed on the ThreeFold Grid. The component itself is designed as reusable infrastructure technology and should be understood by its technical function first, independent of any specific deployment.

## Ownership

This repository is owned and maintained by TF-Tech NV, a Belgian company responsible for the development and maintenance of this technology.

## Requirements

- [`Pulumi`](https://www.pulumi.com/docs/install/) >= 3.128.0
- [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
- [`Go`](https://golang.org/doc/install) >= 1.22

## Using the provider

- Install the provider

```bash
make install_latest
```

### Running the examples

To run examples, make sure you have a mnemonic and a network set.

```bash
export MNEMONIC="mnemonic words"
export NETWORK="network" # dev, qa, test, main -> default is dev
```

- Go to the examples directory `cd examples/go/virtual_machine`
- Run the example `make run`
- To cleanup the resources that you created `make destroy`

## Building The Provider (for development only)

```bash
make build
```

## Run tests

```bash
export MNEMONIC="mnemonic words"
export NETWORK="network" # dev, qa, test, main -> default is dev
```

### Integration tests

```bash
make integration
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
