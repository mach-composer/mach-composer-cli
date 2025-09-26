package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
)

type ApplyOptions struct {
	ForceInit             bool
	IgnoreChangeDetection bool
	BufferLogs            bool
	Github                bool
	Destroy               bool
	AutoApprove           bool
	Filters               []string
}

type PlanOptions struct {
	ForceInit             bool
	IgnoreChangeDetection bool
	BufferLogs            bool
	Github                bool
	Lock                  bool
	Filters               []string
}

type ProxyOptions struct {
	IgnoreChangeDetection bool
	BufferLogs            bool
	Github                bool
	Command               []string
	Filters               []string
}

type ShowPlanOptions struct {
	ForceInit             bool
	IgnoreChangeDetection bool
	BufferLogs            bool
	Github                bool
	NoColor               bool
	Filters               []string
}

type ValidateOptions struct {
	BufferLogs bool
	Github     bool
	Filters    []string
}

type InitOptions struct {
	BufferLogs bool
	Github     bool
	Filters    []string
}

type Runner interface {
	TerraformApply(ctx context.Context, dg *graph.Graph, opts *ApplyOptions) error
	TerraformInit(ctx context.Context, dg *graph.Graph, opts *InitOptions) error
	TerraformPlan(ctx context.Context, dg *graph.Graph, opts *PlanOptions) error
	TerraformProxy(ctx context.Context, dg *graph.Graph, opts *ProxyOptions) error
	TerraformShow(ctx context.Context, dg *graph.Graph, opts *ShowPlanOptions) error
}
