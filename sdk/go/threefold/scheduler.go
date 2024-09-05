// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package threefold

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/threefoldtech/pulumi-threefold/sdk/go/threefold/internal"
)

type Scheduler struct {
	pulumi.CustomResourceState

	Certification_type pulumi.StringPtrOutput `pulumi:"certification_type"`
	City               pulumi.StringPtrOutput `pulumi:"city"`
	Country            pulumi.StringPtrOutput `pulumi:"country"`
	Dedicated          pulumi.BoolPtrOutput   `pulumi:"dedicated"`
	Domain             pulumi.BoolPtrOutput   `pulumi:"domain"`
	Farm_ids           pulumi.IntArrayOutput  `pulumi:"farm_ids"`
	Farm_name          pulumi.StringPtrOutput `pulumi:"farm_name"`
	Free_ips           pulumi.IntPtrOutput    `pulumi:"free_ips"`
	Gpu_available      pulumi.BoolPtrOutput   `pulumi:"gpu_available"`
	Gpu_device_id      pulumi.StringPtrOutput `pulumi:"gpu_device_id"`
	Gpu_device_name    pulumi.StringPtrOutput `pulumi:"gpu_device_name"`
	Gpu_vendor_id      pulumi.StringPtrOutput `pulumi:"gpu_vendor_id"`
	Gpu_vendor_name    pulumi.StringPtrOutput `pulumi:"gpu_vendor_name"`
	Has_gpu            pulumi.BoolPtrOutput   `pulumi:"has_gpu"`
	Hru                pulumi.IntPtrOutput    `pulumi:"hru"`
	Ipv4               pulumi.BoolPtrOutput   `pulumi:"ipv4"`
	Ipv6               pulumi.BoolPtrOutput   `pulumi:"ipv6"`
	Mru                pulumi.IntPtrOutput    `pulumi:"mru"`
	Node_id            pulumi.IntPtrOutput    `pulumi:"node_id"`
	Nodes              pulumi.IntArrayOutput  `pulumi:"nodes"`
	Rentable           pulumi.BoolPtrOutput   `pulumi:"rentable"`
	Rented             pulumi.BoolPtrOutput   `pulumi:"rented"`
	Sru                pulumi.IntPtrOutput    `pulumi:"sru"`
	Twin_id            pulumi.IntPtrOutput    `pulumi:"twin_id"`
}

// NewScheduler registers a new resource with the given unique name, arguments, and options.
func NewScheduler(ctx *pulumi.Context,
	name string, args *SchedulerArgs, opts ...pulumi.ResourceOption) (*Scheduler, error) {
	if args == nil {
		args = &SchedulerArgs{}
	}

	opts = internal.PkgResourceDefaultOpts(opts)
	var resource Scheduler
	err := ctx.RegisterResource("threefold:index:Scheduler", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetScheduler gets an existing Scheduler resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetScheduler(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *SchedulerState, opts ...pulumi.ResourceOption) (*Scheduler, error) {
	var resource Scheduler
	err := ctx.ReadResource("threefold:index:Scheduler", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering Scheduler resources.
type schedulerState struct {
}

type SchedulerState struct {
}

func (SchedulerState) ElementType() reflect.Type {
	return reflect.TypeOf((*schedulerState)(nil)).Elem()
}

type schedulerArgs struct {
	Certification_type *string `pulumi:"certification_type"`
	City               *string `pulumi:"city"`
	Country            *string `pulumi:"country"`
	Dedicated          *bool   `pulumi:"dedicated"`
	Domain             *bool   `pulumi:"domain"`
	Farm_ids           []int   `pulumi:"farm_ids"`
	Farm_name          *string `pulumi:"farm_name"`
	Free_ips           *int    `pulumi:"free_ips"`
	Gpu_available      *bool   `pulumi:"gpu_available"`
	Gpu_device_id      *string `pulumi:"gpu_device_id"`
	Gpu_device_name    *string `pulumi:"gpu_device_name"`
	Gpu_vendor_id      *string `pulumi:"gpu_vendor_id"`
	Gpu_vendor_name    *string `pulumi:"gpu_vendor_name"`
	Has_gpu            *bool   `pulumi:"has_gpu"`
	Hru                *int    `pulumi:"hru"`
	Ipv4               *bool   `pulumi:"ipv4"`
	Ipv6               *bool   `pulumi:"ipv6"`
	Mru                *int    `pulumi:"mru"`
	Node_id            *int    `pulumi:"node_id"`
	Rentable           *bool   `pulumi:"rentable"`
	Rented             *bool   `pulumi:"rented"`
	Sru                *int    `pulumi:"sru"`
	Twin_id            *int    `pulumi:"twin_id"`
}

// The set of arguments for constructing a Scheduler resource.
type SchedulerArgs struct {
	Certification_type pulumi.StringPtrInput
	City               pulumi.StringPtrInput
	Country            pulumi.StringPtrInput
	Dedicated          pulumi.BoolPtrInput
	Domain             pulumi.BoolPtrInput
	Farm_ids           pulumi.IntArrayInput
	Farm_name          pulumi.StringPtrInput
	Free_ips           pulumi.IntPtrInput
	Gpu_available      pulumi.BoolPtrInput
	Gpu_device_id      pulumi.StringPtrInput
	Gpu_device_name    pulumi.StringPtrInput
	Gpu_vendor_id      pulumi.StringPtrInput
	Gpu_vendor_name    pulumi.StringPtrInput
	Has_gpu            pulumi.BoolPtrInput
	Hru                pulumi.IntPtrInput
	Ipv4               pulumi.BoolPtrInput
	Ipv6               pulumi.BoolPtrInput
	Mru                pulumi.IntPtrInput
	Node_id            pulumi.IntPtrInput
	Rentable           pulumi.BoolPtrInput
	Rented             pulumi.BoolPtrInput
	Sru                pulumi.IntPtrInput
	Twin_id            pulumi.IntPtrInput
}

func (SchedulerArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*schedulerArgs)(nil)).Elem()
}

type SchedulerInput interface {
	pulumi.Input

	ToSchedulerOutput() SchedulerOutput
	ToSchedulerOutputWithContext(ctx context.Context) SchedulerOutput
}

func (*Scheduler) ElementType() reflect.Type {
	return reflect.TypeOf((**Scheduler)(nil)).Elem()
}

func (i *Scheduler) ToSchedulerOutput() SchedulerOutput {
	return i.ToSchedulerOutputWithContext(context.Background())
}

func (i *Scheduler) ToSchedulerOutputWithContext(ctx context.Context) SchedulerOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SchedulerOutput)
}

// SchedulerArrayInput is an input type that accepts SchedulerArray and SchedulerArrayOutput values.
// You can construct a concrete instance of `SchedulerArrayInput` via:
//
//	SchedulerArray{ SchedulerArgs{...} }
type SchedulerArrayInput interface {
	pulumi.Input

	ToSchedulerArrayOutput() SchedulerArrayOutput
	ToSchedulerArrayOutputWithContext(context.Context) SchedulerArrayOutput
}

type SchedulerArray []SchedulerInput

func (SchedulerArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*Scheduler)(nil)).Elem()
}

func (i SchedulerArray) ToSchedulerArrayOutput() SchedulerArrayOutput {
	return i.ToSchedulerArrayOutputWithContext(context.Background())
}

func (i SchedulerArray) ToSchedulerArrayOutputWithContext(ctx context.Context) SchedulerArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SchedulerArrayOutput)
}

// SchedulerMapInput is an input type that accepts SchedulerMap and SchedulerMapOutput values.
// You can construct a concrete instance of `SchedulerMapInput` via:
//
//	SchedulerMap{ "key": SchedulerArgs{...} }
type SchedulerMapInput interface {
	pulumi.Input

	ToSchedulerMapOutput() SchedulerMapOutput
	ToSchedulerMapOutputWithContext(context.Context) SchedulerMapOutput
}

type SchedulerMap map[string]SchedulerInput

func (SchedulerMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*Scheduler)(nil)).Elem()
}

func (i SchedulerMap) ToSchedulerMapOutput() SchedulerMapOutput {
	return i.ToSchedulerMapOutputWithContext(context.Background())
}

func (i SchedulerMap) ToSchedulerMapOutputWithContext(ctx context.Context) SchedulerMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(SchedulerMapOutput)
}

type SchedulerOutput struct{ *pulumi.OutputState }

func (SchedulerOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Scheduler)(nil)).Elem()
}

func (o SchedulerOutput) ToSchedulerOutput() SchedulerOutput {
	return o
}

func (o SchedulerOutput) ToSchedulerOutputWithContext(ctx context.Context) SchedulerOutput {
	return o
}

func (o SchedulerOutput) Certification_type() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.StringPtrOutput { return v.Certification_type }).(pulumi.StringPtrOutput)
}

func (o SchedulerOutput) City() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.StringPtrOutput { return v.City }).(pulumi.StringPtrOutput)
}

func (o SchedulerOutput) Country() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.StringPtrOutput { return v.Country }).(pulumi.StringPtrOutput)
}

func (o SchedulerOutput) Dedicated() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.BoolPtrOutput { return v.Dedicated }).(pulumi.BoolPtrOutput)
}

func (o SchedulerOutput) Domain() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.BoolPtrOutput { return v.Domain }).(pulumi.BoolPtrOutput)
}

func (o SchedulerOutput) Farm_ids() pulumi.IntArrayOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.IntArrayOutput { return v.Farm_ids }).(pulumi.IntArrayOutput)
}

func (o SchedulerOutput) Farm_name() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.StringPtrOutput { return v.Farm_name }).(pulumi.StringPtrOutput)
}

func (o SchedulerOutput) Free_ips() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.IntPtrOutput { return v.Free_ips }).(pulumi.IntPtrOutput)
}

func (o SchedulerOutput) Gpu_available() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.BoolPtrOutput { return v.Gpu_available }).(pulumi.BoolPtrOutput)
}

func (o SchedulerOutput) Gpu_device_id() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.StringPtrOutput { return v.Gpu_device_id }).(pulumi.StringPtrOutput)
}

func (o SchedulerOutput) Gpu_device_name() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.StringPtrOutput { return v.Gpu_device_name }).(pulumi.StringPtrOutput)
}

func (o SchedulerOutput) Gpu_vendor_id() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.StringPtrOutput { return v.Gpu_vendor_id }).(pulumi.StringPtrOutput)
}

func (o SchedulerOutput) Gpu_vendor_name() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.StringPtrOutput { return v.Gpu_vendor_name }).(pulumi.StringPtrOutput)
}

func (o SchedulerOutput) Has_gpu() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.BoolPtrOutput { return v.Has_gpu }).(pulumi.BoolPtrOutput)
}

func (o SchedulerOutput) Hru() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.IntPtrOutput { return v.Hru }).(pulumi.IntPtrOutput)
}

func (o SchedulerOutput) Ipv4() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.BoolPtrOutput { return v.Ipv4 }).(pulumi.BoolPtrOutput)
}

func (o SchedulerOutput) Ipv6() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.BoolPtrOutput { return v.Ipv6 }).(pulumi.BoolPtrOutput)
}

func (o SchedulerOutput) Mru() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.IntPtrOutput { return v.Mru }).(pulumi.IntPtrOutput)
}

func (o SchedulerOutput) Node_id() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.IntPtrOutput { return v.Node_id }).(pulumi.IntPtrOutput)
}

func (o SchedulerOutput) Nodes() pulumi.IntArrayOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.IntArrayOutput { return v.Nodes }).(pulumi.IntArrayOutput)
}

func (o SchedulerOutput) Rentable() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.BoolPtrOutput { return v.Rentable }).(pulumi.BoolPtrOutput)
}

func (o SchedulerOutput) Rented() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.BoolPtrOutput { return v.Rented }).(pulumi.BoolPtrOutput)
}

func (o SchedulerOutput) Sru() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.IntPtrOutput { return v.Sru }).(pulumi.IntPtrOutput)
}

func (o SchedulerOutput) Twin_id() pulumi.IntPtrOutput {
	return o.ApplyT(func(v *Scheduler) pulumi.IntPtrOutput { return v.Twin_id }).(pulumi.IntPtrOutput)
}

type SchedulerArrayOutput struct{ *pulumi.OutputState }

func (SchedulerArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*Scheduler)(nil)).Elem()
}

func (o SchedulerArrayOutput) ToSchedulerArrayOutput() SchedulerArrayOutput {
	return o
}

func (o SchedulerArrayOutput) ToSchedulerArrayOutputWithContext(ctx context.Context) SchedulerArrayOutput {
	return o
}

func (o SchedulerArrayOutput) Index(i pulumi.IntInput) SchedulerOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *Scheduler {
		return vs[0].([]*Scheduler)[vs[1].(int)]
	}).(SchedulerOutput)
}

type SchedulerMapOutput struct{ *pulumi.OutputState }

func (SchedulerMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*Scheduler)(nil)).Elem()
}

func (o SchedulerMapOutput) ToSchedulerMapOutput() SchedulerMapOutput {
	return o
}

func (o SchedulerMapOutput) ToSchedulerMapOutputWithContext(ctx context.Context) SchedulerMapOutput {
	return o
}

func (o SchedulerMapOutput) MapIndex(k pulumi.StringInput) SchedulerOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *Scheduler {
		return vs[0].(map[string]*Scheduler)[vs[1].(string)]
	}).(SchedulerOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*SchedulerInput)(nil)).Elem(), &Scheduler{})
	pulumi.RegisterInputType(reflect.TypeOf((*SchedulerArrayInput)(nil)).Elem(), SchedulerArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*SchedulerMapInput)(nil)).Elem(), SchedulerMap{})
	pulumi.RegisterOutputType(SchedulerOutput{})
	pulumi.RegisterOutputType(SchedulerArrayOutput{})
	pulumi.RegisterOutputType(SchedulerMapOutput{})
}