package cmd

import (
	"github.com/ReallyLiri/gow/pkg"
	"github.com/spf13/cobra"
	"log"
)

var replacesCmd = &cobra.Command{
	Use:   "replaces",
	Short: "sync 'replace' directives of workspace modules",
	Long:  `An expansion of the 'go work sync' command.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.WorkspaceSyncReplaces()
		if err != nil {
			log.Fatalf("failed to sync workspace replaces: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(replacesCmd)
}
