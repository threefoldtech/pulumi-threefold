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
    public sealed class VMInput
    {
        public readonly int Cpu;
        public readonly string? Description;
        public readonly string? Entrypoint;
        public readonly ImmutableDictionary<string, string>? Env_vars;
        public readonly string Flist;
        public readonly string? Flist_checksum;
        public readonly ImmutableArray<string> Gpus;
        public readonly int Memory;
        public readonly ImmutableArray<Outputs.Mount> Mounts;
        public readonly string Name;
        public readonly string Network_name;
        public readonly bool? Planetary;
        public readonly bool? Public_ip;
        public readonly bool? Public_ip6;
        public readonly int? Rootfs_size;
        public readonly ImmutableArray<Outputs.Zlog> Zlogs;

        [OutputConstructor]
        private VMInput(
            int cpu,

            string? description,

            string? entrypoint,

            ImmutableDictionary<string, string>? env_vars,

            string flist,

            string? flist_checksum,

            ImmutableArray<string> gpus,

            int memory,

            ImmutableArray<Outputs.Mount> mounts,

            string name,

            string network_name,

            bool? planetary,

            bool? public_ip,

            bool? public_ip6,

            int? rootfs_size,

            ImmutableArray<Outputs.Zlog> zlogs)
        {
            Cpu = cpu;
            Description = description;
            Entrypoint = entrypoint;
            Env_vars = env_vars;
            Flist = flist;
            Flist_checksum = flist_checksum;
            Gpus = gpus;
            Memory = memory;
            Mounts = mounts;
            Name = name;
            Network_name = network_name;
            Planetary = planetary;
            Public_ip = public_ip;
            Public_ip6 = public_ip6;
            Rootfs_size = rootfs_size;
            Zlogs = zlogs;
        }
    }
}