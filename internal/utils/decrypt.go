package utils

import (
	"context"
	"os"
)

// decryptYaml takes a filename and returns the decrypted yaml.
// This command directly calls the sops binary instead of using the
// go.mozilla.org/sops/v3/decrypt package since that adds numerous dependencies
// and adds ~19mb to the generated binary
func DecryptYaml(ctx context.Context, filename string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return RunSops(ctx, wd, "-d", filename, "--output-type=yaml")
}
