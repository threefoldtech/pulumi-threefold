// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Grid.Provider.Outputs
{

    [OutputType]
    public sealed class Group
    {
        public readonly ImmutableArray<Outputs.Backend> Backends;

        [OutputConstructor]
        private Group(ImmutableArray<Outputs.Backend> backends)
        {
            Backends = backends;
        }
    }
}