package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type GroupedError struct {
	msg    string
	Errors []error
}

func NewGroupedError(msg string, errors []error) *GroupedError {
	return &GroupedError{msg: msg, Errors: errors}
}

func (b *GroupedError) Error() string {
	return b.msg
}

type DeprecationOptions struct {
	Site      string
	Component string
	Message   string
	Details   string
	Version   string
}

func DeprecationWarning(opts *DeprecationOptions) {
	if opts.Details != "" {
		opts.Details = utils.TrimIndent(opts.Details)
	}

	log.Warn().
		Str("site", opts.Site).
		Str("component", opts.Component).
		Str("deprecation", opts.Version).
		Str("details", opts.Details).
		Msg(opts.Message)
}

func PrintExitError(summary string, detail ...string) {
	log.Error().
		Str("details", strings.Join(detail, "\n")).
		Msg(summary)
	os.Exit(1)
}

func HandleErr(err error) {
	log.Error().Msgf("Error: %v\n", err)
	var openApiErr *mccsdk.GenericOpenAPIError
	if errors.As(err, &openApiErr) {
		remoteErr := openApiErr.Model()

		switch svcErr := remoteErr.(type) {
		case mccsdk.ErrorUnauthorized:
			PrintExitError("Not authorized", svcErr.GetMessage())
		case mccsdk.ErrorForbidden:
			if svcErr.HasStatus() {
				PrintExitError(svcErr.GetSummary(), svcErr.GetDescription())
			}
			PrintExitError(svcErr.GetMessage(), "Invalid token. Did you run `mach-composer cloud login`?")

		case mccsdk.Error:
			PrintExitError(svcErr.GetSummary(), pie.Map(svcErr.Errors, func(e mccsdk.ErrorObject) string {
				return e.Message
			})...)

		default:
			PrintExitError("Internal error:", err.Error())
		}
		return
	}

	var groupedErr *GroupedError
	if errors.As(err, &groupedErr) {
		var details []string
		for _, e := range groupedErr.Errors {
			details = append(details, e.Error())
		}

		PrintExitError(fmt.Sprintf("A grouped error occured: %s", groupedErr.Error()), details...)
	}

	PrintExitError("An error occurred:", err.Error())
}
