package config

import (
	"context"
	"os"

	"github.com/labd/mach-composer/internal/utils"
)

// decryptYaml takes a filename and returns the decrypted yaml.
// This command directly calls the sops binary instead of using the
// go.mozilla.org/sops/v3/decrypt package since that adds numerous dependencies
// and adds ~19mb to the generated binary
func decryptYaml(ctx context.Context, filename string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return utils.RunSops(ctx, wd, "-d", filename, "--output-type=yaml")
}
