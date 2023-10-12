# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Callable, Mapping, Optional, Sequence, Union, overload
from .. import _utilities

import types

__config__ = pulumi.Config('threefold')


class _ExportableConfig(types.ModuleType):
    @property
    def key_type(self) -> str:
        """
        The key type registered on substrate (ed25519 or sr25519).
        """
        return __config__.get('key_type') or (_utilities.get_env('') or 'sr25519')

    @property
    def mnemonic(self) -> str:
        """
        The mnemonic of the user. It is very secret.
        """
        return __config__.get('mnemonic') or (_utilities.get_env('') or '')

    @property
    def network(self) -> str:
        """
        The network to deploy on.
        """
        return __config__.get('network') or (_utilities.get_env('') or '')

    @property
    def relay_url(self) -> Optional[str]:
        """
        The relay url, example: wss://relay.dev.grid.tf
        """
        return __config__.get('relay_url')

    @property
    def rmb_timeout(self) -> Optional[str]:
        """
        The timeout duration in seconds for rmb calls
        """
        return __config__.get('rmb_timeout')

    @property
    def substrate_url(self) -> Optional[str]:
        """
        The substrate url, example: wss://tfchain.dev.grid.tf/ws
        """
        return __config__.get('substrate_url')
