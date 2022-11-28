package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"

	"github.com/labd/mach-composer/internal/config"
)

type GenerateOptions struct {
	OutputPath string
	Site       string
}

func WriteFiles(cfg *config.MachConfig, options *GenerateOptions) (map[string]string, error) {
	locations := FileLocations(cfg, options)
	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		path := locations[site.Identifier]
		filename := filepath.Join(path, "site.tf")

		fmt.Printf("Generating %s\n", filename)

		body, err := renderSite(cfg, &site)
		if err != nil {
			return nil, err
		}

		// Format and validate the file
		formatted := formatFile([]byte(body))
		if err := validateFile(formatted); err != nil {
			logrus.Error("The generated terraform code is invalid. " +
				"This is a bug in mach composer. Please report the issue at " +
				"https://github.com/labd/mach-composer")
			// os.Exit(255)
		}

		if err := os.MkdirAll(path, 0700); err != nil {
			return nil, fmt.Errorf("error creating directory structure: %w", err)
		}

		if err := os.WriteFile(filename, formatted, 0700); err != nil {
			return nil, fmt.Errorf("error writing file: %w", err)
		}

		// Write extra files
		for extraFilename, content := range cfg.ExtraFiles {
			extraFilename = filepath.Join(path, extraFilename)
			fmt.Printf("Copying %s\n", extraFilename)
			if err := os.WriteFile(extraFilename, content, 0700); err != nil {
				return nil, fmt.Errorf("error writing extra file: %w", err)
			}
		}
	}
	return locations, nil
}

func FileLocations(cfg *config.MachConfig, options *GenerateOptions) map[string]string {
	path := strings.TrimSuffix(filepath.Base(cfg.Filename), filepath.Ext(cfg.Filename))
	sitesPath := filepath.Join(options.OutputPath, path)

	locations := map[string]string{}

	for i := range cfg.Sites {
		site := cfg.Sites[i]
		locations[site.Identifier] = filepath.Join(sitesPath, site.Identifier)
	}
	return locations
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
		logrus.Debugln("Generate HCL has errors:")
		for _, err := range diags.Errs() {
			logrus.Debugln(err)
		}
		return errors.New("generated HCL is invalid")
	}
	return nil
}
