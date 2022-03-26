package runner

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/labd/mach-composer-go/config"
	"github.com/sirupsen/logrus"
)

func TerraformApply(cfg *config.Root, locations map[string]string) {
	ctx := context.Background()

	for i := range cfg.Sites {
		site := cfg.Sites[i]
		TerraformApplySite(ctx, cfg, &site, locations[site.Identifier])
	}
}

func TerraformApplySite(ctx context.Context, cfg *config.Root, site *config.Site, path string) {
	RunTerraform(ctx, path, "init")
}

func RunTerraform(ctx context.Context, cwd string, args ...string) {
	logrus.Debug("Running: terraform %s\n", strings.Join(args, " "))
	cmd := exec.CommandContext(
		ctx,
		"terraform",
		args...,
	)
	cmd.Dir = cwd
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		panic(err)
	}
}
