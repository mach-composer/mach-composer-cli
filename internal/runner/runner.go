package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
)

type ApplyOptions struct {
	ForceInit             bool
	IgnoreChangeDetection bool
	Destroy               bool
	AutoApprove           bool
}

type PlanOptions struct {
	ForceInit             bool
	IgnoreChangeDetection bool
	Lock                  bool
}

type ProxyOptions struct {
	IgnoreChangeDetection bool
	Command               []string
}

type ShowPlanOptions struct {
	IgnoreChangeDetection bool
	NoColor               bool
}

type Runner interface {
	TerraformApply(ctx context.Context, dg *graph.Graph, opts *ApplyOptions) error
	TerraformInit(ctx context.Context, dg *graph.Graph) error
	TerraformPlan(ctx context.Context, dg *graph.Graph, opts *PlanOptions) error
	TerraformProxy(ctx context.Context, dg *graph.Graph, opts *ProxyOptions) error
	TerraformShow(ctx context.Context, dg *graph.Graph, opts *ShowPlanOptions) error
}
