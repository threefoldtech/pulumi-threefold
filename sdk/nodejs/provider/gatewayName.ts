// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "../utilities";

export class GatewayName extends pulumi.CustomResource {
    /**
     * Get an existing GatewayName resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): GatewayName {
        return new GatewayName(name, undefined as any, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'grid:provider:GatewayName';

    /**
     * Returns true if the given object is an instance of GatewayName.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is GatewayName {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === GatewayName.__pulumiType;
    }

    public readonly backends!: pulumi.Output<string[]>;
    public /*out*/ readonly contract_id!: pulumi.Output<number>;
    public readonly description!: pulumi.Output<string | undefined>;
    public /*out*/ readonly fqdn!: pulumi.Output<string>;
    public readonly name!: pulumi.Output<string>;
    public /*out*/ readonly name_contract_id!: pulumi.Output<number>;
    public readonly network!: pulumi.Output<string | undefined>;
    public /*out*/ readonly node_deployment_id!: pulumi.Output<{[key: string]: number}>;
    public readonly node_id!: pulumi.Output<any>;
    public readonly solution_type!: pulumi.Output<string | undefined>;
    public readonly tls_passthrough!: pulumi.Output<boolean | undefined>;

    /**
     * Create a GatewayName resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: GatewayNameArgs, opts?: pulumi.CustomResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            if ((!args || args.backends === undefined) && !opts.urn) {
                throw new Error("Missing required property 'backends'");
            }
            if ((!args || args.name === undefined) && !opts.urn) {
                throw new Error("Missing required property 'name'");
            }
            if ((!args || args.node_id === undefined) && !opts.urn) {
                throw new Error("Missing required property 'node_id'");
            }
            resourceInputs["backends"] = args ? args.backends : undefined;
            resourceInputs["description"] = args ? args.description : undefined;
            resourceInputs["name"] = args ? args.name : undefined;
            resourceInputs["network"] = args ? args.network : undefined;
            resourceInputs["node_id"] = args ? args.node_id : undefined;
            resourceInputs["solution_type"] = args ? args.solution_type : undefined;
            resourceInputs["tls_passthrough"] = args ? args.tls_passthrough : undefined;
            resourceInputs["contract_id"] = undefined /*out*/;
            resourceInputs["fqdn"] = undefined /*out*/;
            resourceInputs["name_contract_id"] = undefined /*out*/;
            resourceInputs["node_deployment_id"] = undefined /*out*/;
        } else {
            resourceInputs["backends"] = undefined /*out*/;
            resourceInputs["contract_id"] = undefined /*out*/;
            resourceInputs["description"] = undefined /*out*/;
            resourceInputs["fqdn"] = undefined /*out*/;
            resourceInputs["name"] = undefined /*out*/;
            resourceInputs["name_contract_id"] = undefined /*out*/;
            resourceInputs["network"] = undefined /*out*/;
            resourceInputs["node_deployment_id"] = undefined /*out*/;
            resourceInputs["node_id"] = undefined /*out*/;
            resourceInputs["solution_type"] = undefined /*out*/;
            resourceInputs["tls_passthrough"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(GatewayName.__pulumiType, name, resourceInputs, opts);
    }
}

/**
 * The set of arguments for constructing a GatewayName resource.
 */
export interface GatewayNameArgs {
    backends: pulumi.Input<pulumi.Input<string>[]>;
    description?: pulumi.Input<string>;
    name: pulumi.Input<string>;
    network?: pulumi.Input<string>;
    node_id: any;
    solution_type?: pulumi.Input<string>;
    tls_passthrough?: pulumi.Input<boolean>;
}