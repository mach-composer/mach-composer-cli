package generator

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/lockfile"
)

type GenerateOptions struct {
	OutputPath string
	Site       string
}

//go:embed templates/*.tmpl
var templates embed.FS

type WriterFunc func(ctx context.Context, cfg *config.MachConfig, g *dependency.Graph, options *GenerateOptions) error

func NewWriterFunc(typ config.DeploymentType) (WriterFunc, error) {
	switch typ {
	case config.Site:
		return sitesWriter, nil
	case config.SiteComponent:
		return componentWriter, nil
	default:
		return nil, fmt.Errorf("unknown deployment type %s", typ)
	}
}

// Write is the main entrypoint for this module. It takes the given MachConfig and graph and iterates the nodes to generate
// the required terraform files.
func Write(ctx context.Context, cfg *config.MachConfig, g *dependency.Graph, options *GenerateOptions) error {
	writer, err := NewWriterFunc(cfg.MachComposer.Deployment.Type)
	if err != nil {
		return err
	}

	return writer(ctx, cfg, g, options)
}

func writeContent(hash, path, content string) error {
	lock, err := lockfile.GetLock(hash, path)
	if err != nil {
		return err
	}

	if !lock.HasChanges(hash) {
		log.Info().Msgf("Files for path %s are up-to-date", path)
		return nil
	}

	filename := filepath.Join(path, "main.tf")

	log.Info().Msgf("Writing %s", filename)

	// Format and validate the file
	formatted := formatFile([]byte(content))
	if err := validateFile(formatted); err != nil {
		log.Error().Msg("The generated terraform code is invalid. " +
			"This is a bug in mach composer. Please report the issue at " +
			"https://github.com/mach-composer/mach-composer-cli")
	}

	if err := os.MkdirAll(path, 0700); err != nil {
		return fmt.Errorf("error creating directory structure: %w", err)
	}

	if err := os.WriteFile(filename, formatted, 0700); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	if err := lock.Update(hash); err != nil {
		return err
	}
	if err := lockfile.WriteLock(lock); err != nil {
		return err
	}

	return nil
}

func FileLocations(_ *config.MachConfig, _ *GenerateOptions) map[string]string {
	panic("implement me")
}

func formatFile(src []byte) []byte {
	// Trim whitespaces prefix
	regex := regexp.MustCompile(`(?m)^\s*`)
	src = regex.ReplaceAll(src, []byte(""))

	// Trim whitespace suffix
	regex = regexp.MustCompile(`(?m)\s*$`)
	src = regex.ReplaceAll(src, []byte(""))

	// Close empty curly blocks on same line
	regex = regexp.MustCompile(`(?m){$\s+}$`)
	src = regex.ReplaceAll(src, []byte("{}"))

	// Close empty array blocks on same line
	regex = regexp.MustCompile(`(?m)\[$\s+\]$`)
	src = regex.ReplaceAll(src, []byte("[]"))

	// Return re-formatted version
	src = hclwrite.Format(src)

	// Insert newline after closing curly brace
	regex = regexp.MustCompile("(?m)^}$")
	src = regex.ReplaceAll(src, []byte("}\n"))

	return src
}

func validateFile(src []byte) error {
	parser := hclparse.NewParser()

	_, diags := parser.ParseHCL(src, "site.tf")
	if diags.HasErrors() {
		log.Debug().Msg("Generate HCL has errors:")
		for _, err := range diags.Errs() {
			log.Debug().Err(err).Msg("error")
		}
		return errors.New("generated HCL is invalid")
	}
	return nil
}

func copyFile(srcPath, dstPath string) error {
	// Open the source file
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Read the contents of the source file
	srcContents, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	// Create the destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// WriteLock the contents of the source file to the destination file
	_, err = dst.Write(srcContents)
	if err != nil {
		return err
	}

	return nil
}
