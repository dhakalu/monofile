package main

import (
	"dhakalu/monofile/internal/builders"
	mr "dhakalu/monofile/internal/monorepo"
	"os"

	"github.com/spf13/cobra"
)

func buildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build",
		Short: "build all projects in the monorepo",
		Run: func(cmd *cobra.Command, args []string) {
			scanner, err := mr.NewScanner("")
			if err != nil {
				cmd.PrintErrf("❌ Error creating scanner: %v", err)
				os.Exit(1)
			}
			scanResult, err := scanner.Scan()
			if err != nil {
				cmd.PrintErrf("❌ Error scanning the monorepo: %v", err)
				os.Exit(1)
			}
			for _, warning := range scanResult.Warnings {
				cmd.Printf(" ⚠️ Skipped %s: %s\n", warning.Path, warning.Message)
			}
			cmd.Printf("Building %d projects\n", len(scanResult.Projects))
			for _, p := range scanResult.Projects {
				cmd.Printf("Building %s", p.Path)
				builder := builders.GetBuilderForLanguage(p.Language)
				if builder == nil {
					cmd.Printf(" ⚠️ No builder found for language %s\n", string(p.Language))
					continue
				}
				info := builder.Build(p)
				if info.Status == builders.BuildStatusFailed {
					cmd.PrintErrf("❌ Error building %s\n Error: %v", p.Path, info.Error)
				}
			}
		},
	}
}
