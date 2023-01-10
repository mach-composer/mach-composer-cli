package cli

import "fmt"

var (
	version = "unknown"
	commit  = "none"
	date    = "unknown"
)

type MetaData struct {
	Version     string
	Commit      string
	ShortCommit string
	Date        string
}

func (md *MetaData) ShortHash() string {
	if len(md.Commit) >= 7 {
		return md.Commit[:7]
	}
	return md.Commit
}

func (md *MetaData) String() string {
	return fmt.Sprintf("mach-composer %s (%s) - %s\n", md.Version, md.ShortCommit, md.Date)
}

func GetVersionMetadata() *MetaData {
	return &MetaData{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}
