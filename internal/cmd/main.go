package main

import (
	"dhakalu/monofile/internal/builders"
	"dhakalu/monofile/internal/model"
	mr "dhakalu/monofile/internal/monorepo"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "monofile",
		Short: "monofile is a simple CLI tool for managing monorepos",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(
				`
monofile is a simple CLI tool for managing monorepos

Usage:
	monofile <command> [arguments]
				
Commands:
	build   Build all projects in the monorepo
			`)
		},
	}

	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "build all projects in the monorepo",
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			slog.Info("Scanning the monorepo", slog.String("root", cwd))
			if err != nil {
				slog.Error("Error getting directory:", slog.Any("error", err))
				os.Exit(1)
			}
			dfs := os.DirFS(cwd)
			md := model.RepoMetadata{
				FileSystem: dfs,
				Root:       cwd,
			}
			scanner := mr.NewScanner(md)
			err = scanner.Scan()
			if err != nil {
				slog.Error("Error scanning the monorepo", slog.Any("error", err))
				os.Exit(1)
			}
			for _, p := range scanner.ProjectsMap {
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
	rootCmd.AddCommand(buildCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Error executing command", slog.Any("error", err))
		os.Exit(1)
	}
}
