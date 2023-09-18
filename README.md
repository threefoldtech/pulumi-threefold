# Pulumi-provider-grid

[![Testing](https://github.com/rawdaGastan/pulumi-provider-grid/actions/workflows/test.yml/badge.svg?branch=development)](https://github.com/rawdaGastan/pulumi-provider-grid/actions/workflows/test.yml) [![Lint](https://github.com/rawdaGastan/pulumi-provider-grid/actions/workflows/lint.yml/badge.svg?branch=development)](https://github.com/rawdaGastan/pulumi-provider-grid/actions/workflows/lint.yml) [![Dependabot](https://badgen.net/badge/Dependabot/enabled/green?icon=dependabot)](https://dependabot.com/)

A pulumi provider for the [threefold grid](https://threefold.io) to manage your infrastructure using pulumi.

## Requirements

- [Pulumi](https://www.pulumi.com/docs/install/) >= 3.81.0
- [Go](https://golang.org/doc/install) >= 1.21

## Using the provider

```bash
cd examples/deployment

export MNEMONICS="mnemonics words"
export NETWORK="network" # dev, qa, test, main

mkdir $PWD/state
pulumi login --cloud-url file://$PWD/state
pulumi stack init test
pulumi up --yes
pulumi up --yes
pulumi destroy --yes
pulumi stack rm --yes
pulumi logout
```

## Building The Provider (for development only)

```bash
make build
```

## Run tests

```bash
export MNEMONICS="mnemonics words"
export NETWORK="network" # dev, qa, test, main
```

- ### Unit tests

  ```bash
  make tests
  ```

- ### Integration tests

  ```bash
  make integration
  ```
