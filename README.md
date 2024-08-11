# Pulumi-threefold

[![Integration tests](https://github.com/threefoldtech/pulumi-threefold/actions/workflows/integration.yaml/badge.svg?branch=development)](https://github.com/threefoldtech/pulumi-threefold/actions/workflows/integration.yaml) [![Lint](https://github.com/threefoldtech/pulumi-threefold/actions/workflows/lint.yaml/badge.svg?branch=development)](https://github.com/threefoldtech/pulumi-threefold/actions/workflows/lint.yaml) [![Dependabot](https://badgen.net/badge/Dependabot/enabled/green?icon=dependabot)](https://dependabot.com/) <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-54%25-brightgreen.svg?longCache=true&style=flat)</a>

A pulumi provider for the [threefold grid](https://threefold.io) to manage your infrastructure using pulumi.
The provider is available at [pulumi registry](https://www.pulumi.com/registry/packages/threefold/).

## Requirements

- [`Pulumi`](https://www.pulumi.com/docs/install/) >= 3.128.0
- [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
- [`Go`](https://golang.org/doc/install) >= 1.21

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
