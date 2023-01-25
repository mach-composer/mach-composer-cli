package cli

import (
	"os"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/rs/zerolog/log"

	"github.com/labd/mach-composer/internal/utils"
)

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
	if openApiErr, ok := err.(*mccsdk.GenericOpenAPIError); ok {
		remoteErr := openApiErr.Model()

		switch svcErr := remoteErr.(type) {
		case mccsdk.ErrorUnauthorized:
			PrintExitError("Not authorized", svcErr.GetMessage())
		case mccsdk.ErrorForbidden:
			PrintExitError(svcErr.GetSummary(), svcErr.GetDescription())
		case mccsdk.Error:
			PrintExitError(svcErr.GetSummary(), pie.Map(svcErr.Errors, func(e mccsdk.ErrorObject) string {
				return e.Message
			})...)

		default:
			PrintExitError("Internal error:", err.Error())
		}
		return
	}

	PrintExitError("An error occured:", err.Error())
}
