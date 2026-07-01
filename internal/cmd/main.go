package main

import (
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
	scan   Scan the monorepo and list all projects
			`)
		},
	}

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan the monorepo and list all projects",
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
			slog.Info("Scan completed. Projects found:", slog.Int("count", len(scanner.ProjectsMap)))
			for path, project := range scanner.ProjectsMap {
				slog.Info("Project found", slog.String("path", path), slog.String("name", project.Name), slog.Any("language", project.Language))
				for _, dep := range project.Dependencies {
					slog.Info("Dependency found", slog.String("name", dep.Name), slog.String("type", string(dep.Type)), slog.String("version", dep.Version))
				}
			}
		},
	}
	rootCmd.AddCommand(scanCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Error executing command", slog.Any("error", err))
		os.Exit(1)
	}
}
