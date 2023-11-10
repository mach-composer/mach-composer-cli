package generator

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func copySecrets(cfg *config.MachConfig, identifier, outPath string) error {
	for _, fs := range cfg.Variables.GetEncryptedSources(identifier) {
		target := filepath.Join(outPath, fs.Filename)
		log.Info().Msgf("Copying %s", target)

		// This can refer to a file outside the current directory, so we need to create the directory structure
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return fmt.Errorf("error creating directory structure for variables: %w", err)
		}

		if err := copyFile(fs.Filename, target); err != nil {
			return fmt.Errorf("error writing extra file: %w", err)
		}
	}

	return nil
}
