package utils

import "fmt"

const (
	IdentifierFieldName = "identifier"
	ArgsFieldName       = "args"
	CommandFieldName    = "command"
	CwdFieldName        = "cwd"
	ModuleName          = "@module"
	NameName            = "name"
	TimestampMain       = "timestamp"
)

func FormatIdentifier(identifier string) string {
	return fmt.Sprintf("[%s]", identifier)
}
