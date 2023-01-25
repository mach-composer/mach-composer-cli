package cloudcmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/cli"
)


func handleError(err error) error {
	if openApiErr, ok := err.(*mccsdk.GenericOpenAPIError); ok {
		remoteErr := openApiErr.Model()
		if svcErr, ok := remoteErr.(mccsdk.Error); ok {
			cli.PrintExitError("Error during API", svcErr.GetError())
		} else {
			var errorMsg string

			errorData := struct {
				Message string `json:"message"`
			}{}
			if err := json.Unmarshal(openApiErr.Body(), &errorData); err == nil {
				errorMsg = errorData.Message
			} else {
				errorMsg = openApiErr.Error()
			}
			fmt.Println("Server returned an error:", errorMsg)
		}
	} else {
		cli.PrintExitError("Server error", err.Error())
	}

	os.Exit(1)
	return nil
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustGetString(cmd *cobra.Command, key string) string {
	value, err := cmd.Flags().GetString(key)
	if err != nil {
		panic(err)
	}
	return value
}

func writeTable(writer io.Writer, header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetBorder(true)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render() // Send output
	fmt.Println()
}
