package cloudcmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"io"
)

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

func writeTable(writer io.Writer, header []string, data [][]string) error {
	table := tablewriter.NewWriter(writer)
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.Row.Formatting.AutoWrap = tw.WrapNone
		cfg.Header.Formatting.AutoFormat = tw.On
		cfg.Header.Alignment.Global = tw.AlignLeft
		cfg.Row.Alignment.Global = tw.AlignLeft
		cfg.Behavior.TrimSpace = tw.On
	})
	table.Header(header)
	_ = table.Bulk(data)

	return table.Render()
}
