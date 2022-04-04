package updater

import (
	"fmt"
	"strings"

	"github.com/labd/mach-composer-go/config"
)

type UpdateSet struct {
	filename string
	updates  []ChangeSet
}

type ChangeSet struct {
	LastVersion string
	Changes     []gitCommit
	Component   *config.Component
}

func (cs *ChangeSet) HasChanges() bool {
	return cs.Component.Version != cs.LastVersion
}

func OutputChanges(cs *ChangeSet) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Updates for %s\n", cs.Component.Name)

	if !cs.HasChanges() {
		fmt.Fprintln(&b, "  No updates...")
		fmt.Fprintln(&b, "")
		return b.String()
	}

	for _, commit := range cs.Changes {
		fmt.Fprintf(&b, "  %s: %s <%s>\n", commit.Commit, commit.Message, commit.Author)
	}
	fmt.Fprintln(&b, "")

	return b.String()
}

func (u *UpdateSet) ChangeLog() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Updated %d components\n\n", len(u.updates))
	for _, cs := range u.updates {
		content := OutputChanges(&cs)
		b.WriteString(content)
	}

	return b.String()

}

func (u *UpdateSet) HasChanges() bool {
	for _, cs := range u.updates {
		if cs.HasChanges() {
			return true
		}
	}
	return false
}
