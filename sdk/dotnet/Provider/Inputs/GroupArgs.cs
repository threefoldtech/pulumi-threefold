// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Grid.Provider.Inputs
{

    public sealed class GroupArgs : global::Pulumi.ResourceArgs
    {
        [Input("backends")]
        private InputList<Inputs.BackendArgs>? _backends;
        public InputList<Inputs.BackendArgs> Backends
        {
            get => _backends ?? (_backends = new InputList<Inputs.BackendArgs>());
            set => _backends = value;
        }

        public GroupArgs()
        {
        }
        public static new GroupArgs Empty => new GroupArgs();
    }
}