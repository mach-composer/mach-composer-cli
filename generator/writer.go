package generator

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/labd/mach-composer-go/config"
	"github.com/sirupsen/logrus"
)

func WriteFiles(cfg *config.MachConfig, target string) (map[string]string, error) {

	path := strings.TrimSuffix(filepath.Base(cfg.Filename), filepath.Ext(cfg.Filename))
	sitesPath := filepath.Join(target, path)

	locations := map[string]string{}

	for i := range cfg.Sites {
		site := cfg.Sites[i]

		filename := filepath.Join(sitesPath, site.Identifier, "site.tf")
		logrus.Infof("Generating %s\n", filename)

		body, err := Render(cfg, &site)
		if err != nil {
			panic(err)
		}

		// Format and validate the file
		formatted := FormatFile([]byte(body))
		if err := ValidateFile(formatted); err != nil {
			logrus.Error("The generated terraform code is invalid. " +
				"This is a bug in mach composer. Please report the issue at " +
				"https://github.com/labd/mach-composer")
			// os.Exit(255)
		}

		if err := os.MkdirAll(filepath.Join(sitesPath, site.Identifier), 0700); err != nil {
			panic(err)
		}

		if err := os.WriteFile(filename, formatted, 0700); err != nil {
			panic(err)
		}

		locations[site.Identifier] = filepath.Dir(filename)
	}
	return locations, nil
}

func FormatFile(src []byte) []byte {
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

func ValidateFile(src []byte) error {
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
