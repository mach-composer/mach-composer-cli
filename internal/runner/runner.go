package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
)

type ApplyOptions struct {
	Force       bool
	Destroy     bool
	AutoApprove bool
}

type PlanOptions struct {
	Force bool
	Lock  bool
}

type ProxyOptions struct {
	Force   bool
	Command []string
}

type ShowPlanOptions struct {
	Force   bool
	NoColor bool
}

type Runner interface {
	TerraformApply(ctx context.Context, dg *graph.Graph, opts *ApplyOptions) error
	TerraformInit(ctx context.Context, dg *graph.Graph) error
	TerraformPlan(ctx context.Context, dg *graph.Graph, opts *PlanOptions) error
	TerraformProxy(ctx context.Context, dg *graph.Graph, opts *ProxyOptions) error
	TerraformShow(ctx context.Context, dg *graph.Graph, opts *ShowPlanOptions) error
}
