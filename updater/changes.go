package updater

import (
	"fmt"

	"github.com/labd/mach-composer-go/config"
)

type ChangeSet struct {
	LastVersion string
	Changes     []gitCommit
	Component   *config.Component
}

func (cs *ChangeSet) HasChanges() bool {
	return cs.Component.Version != cs.LastVersion
}

func OutputChanges(cs *ChangeSet) {
	fmt.Printf("Updates for %s...\n", cs.Component.Name)

	if !cs.HasChanges() {
		fmt.Println("  No updates...")
		fmt.Println("")
		return
	}

	for _, commit := range cs.Changes {
		fmt.Printf("  %s: %s <%s>\n", commit.Commit, commit.Message, commit.Author)
	}
	fmt.Println("")
}
