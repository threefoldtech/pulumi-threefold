# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from .. import _utilities
from . import outputs

__all__ = [
    'Backend',
    'Disk',
    'Group',
    'K8sNodeComputed',
    'K8sNodeInput',
    'Metadata',
    'Mount',
    'QSFSComputed',
    'QSFSInput',
    'VMComputed',
    'VMInput',
    'ZDBComputed',
    'ZDBInput',
    'Zlog',
]

@pulumi.output_type
class Backend(dict):
    def __init__(__self__, *,
                 address: str,
                 namespace: str,
                 password: str):
        pulumi.set(__self__, "address", address)
        pulumi.set(__self__, "namespace", namespace)
        pulumi.set(__self__, "password", password)

    @property
    @pulumi.getter
    def address(self) -> str:
        return pulumi.get(self, "address")

    @property
    @pulumi.getter
    def namespace(self) -> str:
        return pulumi.get(self, "namespace")

    @property
    @pulumi.getter
    def password(self) -> str:
        return pulumi.get(self, "password")


@pulumi.output_type
class Disk(dict):
    def __init__(__self__, *,
                 name: str,
                 size: int,
                 description: Optional[str] = None):
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "size", size)
        if description is not None:
            pulumi.set(__self__, "description", description)

    @property
    @pulumi.getter
    def name(self) -> str:
        return pulumi.get(self, "name")

    @property
    @pulumi.getter
    def size(self) -> int:
        return pulumi.get(self, "size")

    @property
    @pulumi.getter
    def description(self) -> Optional[str]:
        return pulumi.get(self, "description")


@pulumi.output_type
class Group(dict):
    def __init__(__self__, *,
                 backends: Optional[Sequence['outputs.Backend']] = None):
        if backends is not None:
            pulumi.set(__self__, "backends", backends)

    @property
    @pulumi.getter
    def backends(self) -> Optional[Sequence['outputs.Backend']]:
        return pulumi.get(self, "backends")


@pulumi.output_type
class K8sNodeComputed(dict):
    def __init__(__self__, *,
                 computed_ip: str,
                 computed_ip6: str,
                 console_url: str,
                 ip: str,
                 network_name: str,
                 ssh_key: str,
                 token: str,
                 ygg_ip: str):
        pulumi.set(__self__, "computed_ip", computed_ip)
        pulumi.set(__self__, "computed_ip6", computed_ip6)
        pulumi.set(__self__, "console_url", console_url)
        pulumi.set(__self__, "ip", ip)
        pulumi.set(__self__, "network_name", network_name)
        pulumi.set(__self__, "ssh_key", ssh_key)
        pulumi.set(__self__, "token", token)
        pulumi.set(__self__, "ygg_ip", ygg_ip)

    @property
    @pulumi.getter
    def computed_ip(self) -> str:
        return pulumi.get(self, "computed_ip")

    @property
    @pulumi.getter
    def computed_ip6(self) -> str:
        return pulumi.get(self, "computed_ip6")

    @property
    @pulumi.getter
    def console_url(self) -> str:
        return pulumi.get(self, "console_url")

    @property
    @pulumi.getter
    def ip(self) -> str:
        return pulumi.get(self, "ip")

    @property
    @pulumi.getter
    def network_name(self) -> str:
        return pulumi.get(self, "network_name")

    @property
    @pulumi.getter
    def ssh_key(self) -> str:
        return pulumi.get(self, "ssh_key")

    @property
    @pulumi.getter
    def token(self) -> str:
        return pulumi.get(self, "token")

    @property
    @pulumi.getter
    def ygg_ip(self) -> str:
        return pulumi.get(self, "ygg_ip")


@pulumi.output_type
class K8sNodeInput(dict):
    def __init__(__self__, *,
                 cpu: int,
                 disk_size: int,
                 memory: int,
                 name: str,
                 node: Any,
                 flist: Optional[str] = None,
                 flist_checksum: Optional[str] = None,
                 planetary: Optional[bool] = None,
                 public_ip: Optional[bool] = None,
                 public_ip6: Optional[bool] = None):
        pulumi.set(__self__, "cpu", cpu)
        pulumi.set(__self__, "disk_size", disk_size)
        pulumi.set(__self__, "memory", memory)
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "node", node)
        if flist is not None:
            pulumi.set(__self__, "flist", flist)
        if flist_checksum is not None:
            pulumi.set(__self__, "flist_checksum", flist_checksum)
        if planetary is not None:
            pulumi.set(__self__, "planetary", planetary)
        if public_ip is not None:
            pulumi.set(__self__, "public_ip", public_ip)
        if public_ip6 is not None:
            pulumi.set(__self__, "public_ip6", public_ip6)

    @property
    @pulumi.getter
    def cpu(self) -> int:
        return pulumi.get(self, "cpu")

    @property
    @pulumi.getter
    def disk_size(self) -> int:
        return pulumi.get(self, "disk_size")

    @property
    @pulumi.getter
    def memory(self) -> int:
        return pulumi.get(self, "memory")

    @property
    @pulumi.getter
    def name(self) -> str:
        return pulumi.get(self, "name")

    @property
    @pulumi.getter
    def node(self) -> Any:
        return pulumi.get(self, "node")

    @property
    @pulumi.getter
    def flist(self) -> Optional[str]:
        return pulumi.get(self, "flist")

    @property
    @pulumi.getter
    def flist_checksum(self) -> Optional[str]:
        return pulumi.get(self, "flist_checksum")

    @property
    @pulumi.getter
    def planetary(self) -> Optional[bool]:
        return pulumi.get(self, "planetary")

    @property
    @pulumi.getter
    def public_ip(self) -> Optional[bool]:
        return pulumi.get(self, "public_ip")

    @property
    @pulumi.getter
    def public_ip6(self) -> Optional[bool]:
        return pulumi.get(self, "public_ip6")


@pulumi.output_type
class Metadata(dict):
    def __init__(__self__, *,
                 encryption_key: str,
                 prefix: str,
                 backends: Optional[Sequence['outputs.Backend']] = None,
                 encryption_algorithm: Optional[str] = None,
                 type: Optional[str] = None):
        pulumi.set(__self__, "encryption_key", encryption_key)
        pulumi.set(__self__, "prefix", prefix)
        if backends is not None:
            pulumi.set(__self__, "backends", backends)
        if encryption_algorithm is not None:
            pulumi.set(__self__, "encryption_algorithm", encryption_algorithm)
        if type is not None:
            pulumi.set(__self__, "type", type)

    @property
    @pulumi.getter
    def encryption_key(self) -> str:
        return pulumi.get(self, "encryption_key")

    @property
    @pulumi.getter
    def prefix(self) -> str:
        return pulumi.get(self, "prefix")

    @property
    @pulumi.getter
    def backends(self) -> Optional[Sequence['outputs.Backend']]:
        return pulumi.get(self, "backends")

    @property
    @pulumi.getter
    def encryption_algorithm(self) -> Optional[str]:
        return pulumi.get(self, "encryption_algorithm")

    @property
    @pulumi.getter
    def type(self) -> Optional[str]:
        return pulumi.get(self, "type")


@pulumi.output_type
class Mount(dict):
    def __init__(__self__, *,
                 disk_name: str,
                 mount_point: str):
        pulumi.set(__self__, "disk_name", disk_name)
        pulumi.set(__self__, "mount_point", mount_point)

    @property
    @pulumi.getter
    def disk_name(self) -> str:
        return pulumi.get(self, "disk_name")

    @property
    @pulumi.getter
    def mount_point(self) -> str:
        return pulumi.get(self, "mount_point")


@pulumi.output_type
class QSFSComputed(dict):
    def __init__(__self__, *,
                 metrics_endpoint: str):
        pulumi.set(__self__, "metrics_endpoint", metrics_endpoint)

    @property
    @pulumi.getter
    def metrics_endpoint(self) -> str:
        return pulumi.get(self, "metrics_endpoint")


@pulumi.output_type
class QSFSInput(dict):
    def __init__(__self__, *,
                 cache: int,
                 encryption_key: str,
                 expected_shards: int,
                 groups: Sequence['outputs.Group'],
                 max_zdb_data_dir_size: int,
                 metadata: 'outputs.Metadata',
                 minimal_shards: int,
                 name: str,
                 redundant_groups: int,
                 redundant_nodes: int,
                 compression_algorithm: Optional[str] = None,
                 description: Optional[str] = None,
                 encryption_algorithm: Optional[str] = None):
        pulumi.set(__self__, "cache", cache)
        pulumi.set(__self__, "encryption_key", encryption_key)
        pulumi.set(__self__, "expected_shards", expected_shards)
        pulumi.set(__self__, "groups", groups)
        pulumi.set(__self__, "max_zdb_data_dir_size", max_zdb_data_dir_size)
        pulumi.set(__self__, "metadata", metadata)
        pulumi.set(__self__, "minimal_shards", minimal_shards)
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "redundant_groups", redundant_groups)
        pulumi.set(__self__, "redundant_nodes", redundant_nodes)
        if compression_algorithm is not None:
            pulumi.set(__self__, "compression_algorithm", compression_algorithm)
        if description is not None:
            pulumi.set(__self__, "description", description)
        if encryption_algorithm is not None:
            pulumi.set(__self__, "encryption_algorithm", encryption_algorithm)

    @property
    @pulumi.getter
    def cache(self) -> int:
        return pulumi.get(self, "cache")

    @property
    @pulumi.getter
    def encryption_key(self) -> str:
        return pulumi.get(self, "encryption_key")

    @property
    @pulumi.getter
    def expected_shards(self) -> int:
        return pulumi.get(self, "expected_shards")

    @property
    @pulumi.getter
    def groups(self) -> Sequence['outputs.Group']:
        return pulumi.get(self, "groups")

    @property
    @pulumi.getter
    def max_zdb_data_dir_size(self) -> int:
        return pulumi.get(self, "max_zdb_data_dir_size")

    @property
    @pulumi.getter
    def metadata(self) -> 'outputs.Metadata':
        return pulumi.get(self, "metadata")

    @property
    @pulumi.getter
    def minimal_shards(self) -> int:
        return pulumi.get(self, "minimal_shards")

    @property
    @pulumi.getter
    def name(self) -> str:
        return pulumi.get(self, "name")

    @property
    @pulumi.getter
    def redundant_groups(self) -> int:
        return pulumi.get(self, "redundant_groups")

    @property
    @pulumi.getter
    def redundant_nodes(self) -> int:
        return pulumi.get(self, "redundant_nodes")

    @property
    @pulumi.getter
    def compression_algorithm(self) -> Optional[str]:
        return pulumi.get(self, "compression_algorithm")

    @property
    @pulumi.getter
    def description(self) -> Optional[str]:
        return pulumi.get(self, "description")

    @property
    @pulumi.getter
    def encryption_algorithm(self) -> Optional[str]:
        return pulumi.get(self, "encryption_algorithm")


@pulumi.output_type
class VMComputed(dict):
    def __init__(__self__, *,
                 computed_ip: str,
                 computed_ip6: str,
                 console_url: str,
                 ygg_ip: str,
                 ip: Optional[str] = None):
        pulumi.set(__self__, "computed_ip", computed_ip)
        pulumi.set(__self__, "computed_ip6", computed_ip6)
        pulumi.set(__self__, "console_url", console_url)
        pulumi.set(__self__, "ygg_ip", ygg_ip)
        if ip is not None:
            pulumi.set(__self__, "ip", ip)

    @property
    @pulumi.getter
    def computed_ip(self) -> str:
        return pulumi.get(self, "computed_ip")

    @property
    @pulumi.getter
    def computed_ip6(self) -> str:
        return pulumi.get(self, "computed_ip6")

    @property
    @pulumi.getter
    def console_url(self) -> str:
        return pulumi.get(self, "console_url")

    @property
    @pulumi.getter
    def ygg_ip(self) -> str:
        return pulumi.get(self, "ygg_ip")

    @property
    @pulumi.getter
    def ip(self) -> Optional[str]:
        return pulumi.get(self, "ip")


@pulumi.output_type
class VMInput(dict):
    def __init__(__self__, *,
                 cpu: int,
                 flist: str,
                 memory: int,
                 name: str,
                 network_name: str,
                 description: Optional[str] = None,
                 entrypoint: Optional[str] = None,
                 env_vars: Optional[Mapping[str, str]] = None,
                 flist_checksum: Optional[str] = None,
                 gpus: Optional[Sequence[str]] = None,
                 mounts: Optional[Sequence['outputs.Mount']] = None,
                 planetary: Optional[bool] = None,
                 public_ip: Optional[bool] = None,
                 public_ip6: Optional[bool] = None,
                 rootfs_size: Optional[int] = None,
                 zlogs: Optional[Sequence['outputs.Zlog']] = None):
        pulumi.set(__self__, "cpu", cpu)
        pulumi.set(__self__, "flist", flist)
        pulumi.set(__self__, "memory", memory)
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "network_name", network_name)
        if description is not None:
            pulumi.set(__self__, "description", description)
        if entrypoint is not None:
            pulumi.set(__self__, "entrypoint", entrypoint)
        if env_vars is not None:
            pulumi.set(__self__, "env_vars", env_vars)
        if flist_checksum is not None:
            pulumi.set(__self__, "flist_checksum", flist_checksum)
        if gpus is not None:
            pulumi.set(__self__, "gpus", gpus)
        if mounts is not None:
            pulumi.set(__self__, "mounts", mounts)
        if planetary is not None:
            pulumi.set(__self__, "planetary", planetary)
        if public_ip is not None:
            pulumi.set(__self__, "public_ip", public_ip)
        if public_ip6 is not None:
            pulumi.set(__self__, "public_ip6", public_ip6)
        if rootfs_size is not None:
            pulumi.set(__self__, "rootfs_size", rootfs_size)
        if zlogs is not None:
            pulumi.set(__self__, "zlogs", zlogs)

    @property
    @pulumi.getter
    def cpu(self) -> int:
        return pulumi.get(self, "cpu")

    @property
    @pulumi.getter
    def flist(self) -> str:
        return pulumi.get(self, "flist")

    @property
    @pulumi.getter
    def memory(self) -> int:
        return pulumi.get(self, "memory")

    @property
    @pulumi.getter
    def name(self) -> str:
        return pulumi.get(self, "name")

    @property
    @pulumi.getter
    def network_name(self) -> str:
        return pulumi.get(self, "network_name")

    @property
    @pulumi.getter
    def description(self) -> Optional[str]:
        return pulumi.get(self, "description")

    @property
    @pulumi.getter
    def entrypoint(self) -> Optional[str]:
        return pulumi.get(self, "entrypoint")

    @property
    @pulumi.getter
    def env_vars(self) -> Optional[Mapping[str, str]]:
        return pulumi.get(self, "env_vars")

    @property
    @pulumi.getter
    def flist_checksum(self) -> Optional[str]:
        return pulumi.get(self, "flist_checksum")

    @property
    @pulumi.getter
    def gpus(self) -> Optional[Sequence[str]]:
        return pulumi.get(self, "gpus")

    @property
    @pulumi.getter
    def mounts(self) -> Optional[Sequence['outputs.Mount']]:
        return pulumi.get(self, "mounts")

    @property
    @pulumi.getter
    def planetary(self) -> Optional[bool]:
        return pulumi.get(self, "planetary")

    @property
    @pulumi.getter
    def public_ip(self) -> Optional[bool]:
        return pulumi.get(self, "public_ip")

    @property
    @pulumi.getter
    def public_ip6(self) -> Optional[bool]:
        return pulumi.get(self, "public_ip6")

    @property
    @pulumi.getter
    def rootfs_size(self) -> Optional[int]:
        return pulumi.get(self, "rootfs_size")

    @property
    @pulumi.getter
    def zlogs(self) -> Optional[Sequence['outputs.Zlog']]:
        return pulumi.get(self, "zlogs")


@pulumi.output_type
class ZDBComputed(dict):
    def __init__(__self__, *,
                 ips: Sequence[str],
                 namespace: str,
                 port: int):
        pulumi.set(__self__, "ips", ips)
        pulumi.set(__self__, "namespace", namespace)
        pulumi.set(__self__, "port", port)

    @property
    @pulumi.getter
    def ips(self) -> Sequence[str]:
        return pulumi.get(self, "ips")

    @property
    @pulumi.getter
    def namespace(self) -> str:
        return pulumi.get(self, "namespace")

    @property
    @pulumi.getter
    def port(self) -> int:
        return pulumi.get(self, "port")


@pulumi.output_type
class ZDBInput(dict):
    def __init__(__self__, *,
                 name: str,
                 password: str,
                 size: int,
                 description: Optional[str] = None,
                 mode: Optional[str] = None,
                 public: Optional[bool] = None):
        pulumi.set(__self__, "name", name)
        pulumi.set(__self__, "password", password)
        pulumi.set(__self__, "size", size)
        if description is not None:
            pulumi.set(__self__, "description", description)
        if mode is None:
            mode = (_utilities.get_env('') or 'user')
        if mode is not None:
            pulumi.set(__self__, "mode", mode)
        if public is not None:
            pulumi.set(__self__, "public", public)

    @property
    @pulumi.getter
    def name(self) -> str:
        return pulumi.get(self, "name")

    @property
    @pulumi.getter
    def password(self) -> str:
        return pulumi.get(self, "password")

    @property
    @pulumi.getter
    def size(self) -> int:
        return pulumi.get(self, "size")

    @property
    @pulumi.getter
    def description(self) -> Optional[str]:
        return pulumi.get(self, "description")

    @property
    @pulumi.getter
    def mode(self) -> Optional[str]:
        return pulumi.get(self, "mode")

    @property
    @pulumi.getter
    def public(self) -> Optional[bool]:
        return pulumi.get(self, "public")


@pulumi.output_type
class Zlog(dict):
    def __init__(__self__, *,
                 output: str,
                 zmachine: str):
        pulumi.set(__self__, "output", output)
        pulumi.set(__self__, "zmachine", zmachine)

    @property
    @pulumi.getter
    def output(self) -> str:
        return pulumi.get(self, "output")

    @property
    @pulumi.getter
    def zmachine(self) -> str:
        return pulumi.get(self, "zmachine")

