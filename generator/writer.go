package generator

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/labd/mach-composer-go/config"
)

func WriteFiles(cfg *config.Root, target string) (map[string]string, error) {

	path := strings.TrimSuffix(filepath.Base(cfg.Filename), filepath.Ext(cfg.Filename))
	sitesPath := filepath.Join(target, path)

	locations := map[string]string{}

	for i := range cfg.Sites {
		site := cfg.Sites[i]

		filename := filepath.Join(sitesPath, site.Identifier, "site.tf")
		log.Printf("Generating %s\n", filename)

		body, err := Render(cfg, &site)
		if err != nil {
			panic(err)
		}

		// Validate and format file
		formatted := FormatFile([]byte(body))

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
	parser := hclparse.NewParser()

	// Validate the generated hcl
	_, diags := parser.ParseHCL(src, "site.tf")
	if diags.HasErrors() {
		log.Println("Generate HCL has errors:")
		for _, err := range diags.Errs() {
			log.Println(err)
		}
		return src
	}

	// Return re-formatted version
	return hclwrite.Format(src)
}
