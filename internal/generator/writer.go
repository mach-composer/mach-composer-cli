package generator

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/state"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type GenerateOptions struct{}

//go:embed templates/*.tmpl
var templates embed.FS

// Write is the main entrypoint for this module. It takes the given MachConfig and graph and iterates the nodes to generate
// the required terraform files.
func Write(ctx context.Context, cfg *config.MachConfig, g *graph.Graph, _ *GenerateOptions) error {
	for _, n := range g.Vertices() {
		sr, err := state.NewRenderer(
			state.Type(cfg.Global.TerraformStateProvider),
			n.Identifier(),
			cfg.Global.TerraformConfig.RemoteState,
		)
		if err != nil {
			return err
		}
		err = cfg.StateRepository.Add(sr)
		if err != nil {
			return err
		}

		if site, ok := n.(*graph.Site); ok {
			for _, c := range site.NestedNodes {
				cfg.StateRepository.Alias(n.Identifier(), c.Identifier())
			}
		}
	}

	for _, n := range g.Vertices() {
		switch n.(type) {
		case *graph.Project:
			log.Debug().Msgf("No global files to generate for project %s", n.Path())
			break
		case *graph.Site:
			if err := copySecrets(cfg, n.Identifier(), n.Path()); err != nil {
				return err
			}
			body, err := renderSite(ctx, cfg, n.(*graph.Site))
			if err != nil {
				return err
			}

			if err = writeContent(n.Path(), body); err != nil {
				return err
			}
			break
		case *graph.SiteComponent:
			if err := copySecrets(cfg, n.Identifier(), n.Path()); err != nil {
				return err
			}

			body, err := renderSiteComponent(ctx, cfg, n.(*graph.SiteComponent))
			if err != nil {
				return err
			}

			if err = writeContent(n.Path(), body); err != nil {
				return err
			}
			break
		default:
			return fmt.Errorf("unknown node type %T", n)
		}

	}

	return nil
}

func writeContent(path, content string) error {
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

	log.Info().Msgf("Wrote files for path %s", path)

	return nil
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
