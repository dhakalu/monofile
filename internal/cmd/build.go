package main

import (
	"dhakalu/monofile/internal/builders"
	mr "dhakalu/monofile/internal/monorepo"
	"log/slog"
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
			projectsMap, err := scanner.Scan()
			if err != nil {
				cmd.PrintErrf("❌ Error scanning the monorepo: %v", err)
				os.Exit(1)
			}
			for _, p := range projectsMap {
				builder := builders.GetBuilderForLanguage(p.Language)
				if builder == nil {
					slog.Warn("No builder found for language", slog.String("language", string(p.Language)))
					continue
				}
				info := builder.Build(p)
				if info.Status == builders.BuildStatusFailed {
					slog.Error("Error building ", slog.String("project", p.Path))
				}
			}
		},
	}
}
