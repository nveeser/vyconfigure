package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Prints the current version.",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(appVersion)
		},
	}
}
