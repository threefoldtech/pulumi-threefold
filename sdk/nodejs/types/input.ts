// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../types/input";
import * as outputs from "../types/output";

import * as utilities from "../utilities";

export namespace provider {
    export interface BackendArgs {
        address: pulumi.Input<string>;
        namespace: pulumi.Input<string>;
        password: pulumi.Input<string>;
    }

    export interface DiskArgs {
        description?: pulumi.Input<string>;
        name: pulumi.Input<string>;
        size: pulumi.Input<number>;
    }

    export interface GroupArgs {
        backends?: pulumi.Input<pulumi.Input<inputs.provider.BackendArgs>[]>;
    }

    export interface K8sNodeInputArgs {
        cpu: pulumi.Input<number>;
        disk_size: pulumi.Input<number>;
        flist?: pulumi.Input<string>;
        flist_checksum?: pulumi.Input<string>;
        memory: pulumi.Input<number>;
        name: pulumi.Input<string>;
        node: any;
        planetary?: pulumi.Input<boolean>;
        public_ip?: pulumi.Input<boolean>;
        public_ip6?: pulumi.Input<boolean>;
    }

    export interface MetadataArgs {
        backends?: pulumi.Input<pulumi.Input<inputs.provider.BackendArgs>[]>;
        encryption_algorithm?: pulumi.Input<string>;
        encryption_key: pulumi.Input<string>;
        prefix: pulumi.Input<string>;
        type?: pulumi.Input<string>;
    }

    export interface MountArgs {
        disk_name: pulumi.Input<string>;
        mount_point: pulumi.Input<string>;
    }

    export interface QSFSInputArgs {
        cache: pulumi.Input<number>;
        compression_algorithm?: pulumi.Input<string>;
        description?: pulumi.Input<string>;
        encryption_algorithm?: pulumi.Input<string>;
        encryption_key: pulumi.Input<string>;
        expected_shards: pulumi.Input<number>;
        groups: pulumi.Input<pulumi.Input<inputs.provider.GroupArgs>[]>;
        max_zdb_data_dir_size: pulumi.Input<number>;
        metadata: pulumi.Input<inputs.provider.MetadataArgs>;
        minimal_shards: pulumi.Input<number>;
        name: pulumi.Input<string>;
        redundant_groups: pulumi.Input<number>;
        redundant_nodes: pulumi.Input<number>;
    }

    export interface VMInputArgs {
        cpu: pulumi.Input<number>;
        description?: pulumi.Input<string>;
        entrypoint?: pulumi.Input<string>;
        env_vars?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
        flist: pulumi.Input<string>;
        flist_checksum?: pulumi.Input<string>;
        gpus?: pulumi.Input<pulumi.Input<string>[]>;
        memory: pulumi.Input<number>;
        mounts?: pulumi.Input<pulumi.Input<inputs.provider.MountArgs>[]>;
        name: pulumi.Input<string>;
        network_name: pulumi.Input<string>;
        planetary?: pulumi.Input<boolean>;
        public_ip?: pulumi.Input<boolean>;
        public_ip6?: pulumi.Input<boolean>;
        rootfs_size?: pulumi.Input<number>;
        zlogs?: pulumi.Input<pulumi.Input<inputs.provider.ZlogArgs>[]>;
    }

    export interface ZDBInputArgs {
        description?: pulumi.Input<string>;
        mode?: pulumi.Input<string>;
        name: pulumi.Input<string>;
        password: pulumi.Input<string>;
        public?: pulumi.Input<boolean>;
        size: pulumi.Input<number>;
    }
    /**
     * zdbinputArgsProvideDefaults sets the appropriate defaults for ZDBInputArgs
     */
    export function zdbinputArgsProvideDefaults(val: ZDBInputArgs): ZDBInputArgs {
        return {
            ...val,
            mode: (val.mode) ?? (utilities.getEnv("") || "user"),
        };
    }

    export interface ZlogArgs {
        output: pulumi.Input<string>;
        zmachine: pulumi.Input<string>;
    }
}
