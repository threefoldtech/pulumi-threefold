# Pulumi-provider-grid

[![Testing](https://github.com/threefoldtech/pulumi-provider-grid/actions/workflows/test.yaml/badge.svg?branch=development)](https://github.com/threefoldtech/pulumi-provider-grid/actions/workflows/test.yaml) [![Lint](https://github.com/threefoldtech/pulumi-provider-grid/actions/workflows/lint.yaml/badge.svg?branch=development)](https://github.com/threefoldtech/pulumi-provider-grid/actions/workflows/lint.yaml) [![Dependabot](https://badgen.net/badge/Dependabot/enabled/green?icon=dependabot)](https://dependabot.com/) <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-0%25-brightgreen.svg?longCache=true&style=flat)</a>

A pulumi provider for the [threefold grid](https://threefold.io) to manage your infrastructure using pulumi.

## Requirements

- [Pulumi](https://www.pulumi.com/docs/install/) >= 3.81.0
- [Go](https://golang.org/doc/install) >= 1.21

## Using the provider

- You can try to run examples:

```bash
cd examples/virtual_machine

export MNEMONICS="mnemonics words"
export NETWORK="network" # dev, qa, test, main -> default is dev

make run
```

- to destroy the resources you created:

```bash
make destroy
```

## Building The Provider (for development only)

```bash
make build
```

## Run tests

```bash
export MNEMONICS="mnemonics words"
export NETWORK="network" # dev, qa, test, main -> default is dev
```

- ### Unit tests

  ```bash
  make tests
  ```

- ### Integration tests

  ```bash
  make integration
  ```
