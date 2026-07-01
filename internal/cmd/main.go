package main

import (
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
	check   Check for dependency version conflicts in the monorepo
			`)
		},
	}

	
	rootCmd.AddCommand(buildCmd())
	rootCmd.AddCommand(checkCmd())

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Error executing command", slog.Any("error", err))
		os.Exit(1)
	}
}
