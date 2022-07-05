package updater

import (
	"context"
	"fmt"
	"os"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/utils"
)

// SopsFileWriter updates the contents of a mach file with the updated
// version of the components
func SopsFileWriter(cfg *config.MachConfig, updates *UpdateSet) {

	indexMap := make(map[string]int)

	for i := range cfg.Components {
		indexMap[cfg.Components[i].Name] = i
	}

	ctx := context.Background()

	for _, c := range updates.updates {
		index := indexMap[c.Component.Name]

		result, err := utils.RunSops(ctx, ".",
			"--set",
			fmt.Sprintf(`["components"][%d]["version"] "%s"`, index, c.LastVersion),
			updates.filename,
		)
		if err != nil {
			fmt.Fprint(os.Stderr, result)
		}
	}

}
