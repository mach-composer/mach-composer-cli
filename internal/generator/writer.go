package generator

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/state"
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

// Write is the main entrypoint for this module. It takes the given MachConfig and graph and iterates the nodes to generate
// the required terraform files.
func Write(ctx context.Context, cfg *config.MachConfig, g *dependency.Graph, options *GenerateOptions) error {
	for _, n := range g.Vertices() {
		sr, err := state.NewRenderer(
			state.Type(cfg.Global.TerraformStateProvider),
			n.Identifier(),
			cfg.Global.TerraformConfig.RemoteState,
		)
		if err != nil {
			return err
		}
		err = cfg.StateRepository.Add(sr.Key(), sr)
		if err != nil {
			return err
		}

		if site, ok := n.(*dependency.Site); ok {
			for _, c := range site.NestedSiteComponentConfigs {
				cfg.StateRepository.Alias(n.Identifier(), c.Name)
			}
		}
	}

	for _, n := range g.Vertices() {
		outPath := fmt.Sprintf("%s/%s", options.OutputPath, n.Path())

		switch n.(type) {
		case *dependency.Project:
			log.Debug().Msgf("No global files to generate for project %s", n.Path())
			break
		case *dependency.Site:
			siteConfig := n.(*dependency.Site).SiteConfig
			nestedComponents := n.(*dependency.Site).NestedSiteComponentConfigs

			if err := copySecrets(cfg, siteConfig.Identifier, outPath); err != nil {
				return err
			}
			body, err := renderSite(ctx, cfg, &siteConfig, nestedComponents)
			if err != nil {
				return err
			}

			if err = writeContent(cfg.ConfigHash, outPath, body); err != nil {
				return err
			}

			log.Info().Msgf("Wrote site %s", siteConfig.Identifier)
			break
		case *dependency.SiteComponent:
			s := n.(*dependency.SiteComponent).SiteConfig
			sc := n.(*dependency.SiteComponent).SiteComponentConfig

			if err := copySecrets(cfg, s.Identifier, outPath); err != nil {
				return err
			}

			body, err := renderSiteComponent(ctx, cfg, &s, &sc)
			if err != nil {
				return err
			}

			if err = writeContent(cfg.ConfigHash, outPath, body); err != nil {
				return err
			}

			log.Info().Msgf("Wrote site component %s", n.Path())
			break
		default:
			return fmt.Errorf("unknown node type %T", n)
		}

	}

	return nil
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
