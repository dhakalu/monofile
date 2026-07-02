package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var verbose bool
	var rootCmd = &cobra.Command{
		Use:           "monofile",
		Short:         "monofile is a simple CLI tool for managing monorepos",
		Args:          cobra.NoArgs,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := slog.LevelWarn
			if verbose {
				level = slog.LevelDebug
			}
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: level,
			})))
		},
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose diagnostic logging")
	rootCmd.AddCommand(buildCmd())
	rootCmd.AddCommand(checkCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ %v\n", err)
		os.Exit(1)
	}
}
