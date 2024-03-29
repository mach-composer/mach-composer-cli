package updater

import (
	"fmt"
	"strings"
	"time"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type UpdateSet struct {
	filename string
	updates  []ChangeSet
}

type ChangeSet struct {
	LastVersion string
	Changes     []CommitData
	Component   *config.ComponentConfig
	Forced      bool
}

type CommitData struct {
	Commit    string
	Parents   []string
	Author    CommitAuthor
	Committer CommitAuthor
	Message   string
	Tags      []string
}

type CommitAuthor struct {
	Name  string
	Email string
	Date  time.Time
}

func (cs *ChangeSet) HasChanges() bool {
	return cs.Component.Version != cs.LastVersion
}

func generateTagMessage(tags []string) string {
	if len(tags) == 0 {
		return ""
	}

	tagList := strings.Join(tags, ", ")
	if len(tagList) > 50 {
		return fmt.Sprintf(" (%d tags found) ", len(tags))
	}

	return fmt.Sprintf(" (tags: %s) ", tagList)
}

func OutputChanges(cs *ChangeSet) string {
	var b strings.Builder

	if cs.Forced && len(cs.Changes) == 0 {
		fmt.Fprintf(&b, "Update %s to %s\n", cs.Component.Name, cs.LastVersion)
		return b.String()
	}

	fmt.Fprintf(&b, "Updates for %s (%s...%s)\n", cs.Component.Name, cs.Component.Version, cs.LastVersion)

	if !cs.HasChanges() {
		fmt.Fprintln(&b, "  No updates...")
		fmt.Fprintln(&b, "")
		return b.String()
	}

	for _, commit := range cs.Changes {
		tagMsg := generateTagMessage(commit.Tags)

		fmt.Fprintf(&b, "%s%s: %s (%s <%s> %s)\n",
			commit.Commit,
			tagMsg,
			commit.Message,
			commit.Author.Name,
			commit.Author.Email,
			commit.Author.Date.Format(time.RFC3339),
		)
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

func (u *UpdateSet) ComponentChangeLog(component string) string {
	var b strings.Builder

	for _, cs := range u.updates {
		if strings.EqualFold(cs.Component.Name, component) {
			content := OutputChanges(&cs)
			b.WriteString(content)
		}
	}
	return b.String()
}

func (u *UpdateSet) HasChanges() bool {
	for _, cs := range u.updates {
		if cs.HasChanges() || cs.Forced {
			return true
		}
	}
	return false
}
