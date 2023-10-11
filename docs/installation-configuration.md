---
title: Grid Installation & Configuration
meta_desc: Information on how to install the Grid provider.
layout: package
---

## Installation

The Pulumi Grid provider is available as a package in all Pulumi languages:

* Go: [`github.com/threefold/pulumi-grid/sdk`](https://pkg.go.dev/github.com/threefold/pulumi-grid/sdk)

## Setup

To provision resources with the Pulumi grid provider, you need to provide the `mnemonic`.

## Configuration Options

Use `pulumi config set grid:<option>`.

The following configuration points are available for the `grid` provider:

* `grid:mnemonic` (environment: `MNEMONIC`) -  This is the grid mnemonic. You can create a new account if you don't have [mnemonics](https://threefoldtech.github.io/info_grid/dashboard/portal/dashboard_portal_polkadot_create_account.html).

* `grid:network` (environment: `NETWORK`) - specify which grid network (dev/qa/mainnet/testnet) to deploy on (default is dev).
