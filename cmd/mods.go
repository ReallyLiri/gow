package cmd

import (
	"github.com/ReallyLiri/gow/pkg"
	"github.com/spf13/cobra"
	"log"
)

var modsCmd = &cobra.Command{
	Use:   "mods",
	Short: "Run a command for each mod",
	Long:  `Run a command for each mod. For example: ...`,
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.RunForEachMod(args)
		if err != nil {
			log.Fatalf("failed to run command for each mod: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(modsCmd)
}
