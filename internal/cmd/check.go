package main

import (
	mr "dhakalu/monofile/internal/monorepo"
	"os"

	"github.com/spf13/cobra"
)

func checkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "check for dependency version conflicts in the monorepo",
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
			versions := mr.GetDependencyVersions(projectsMap)
			for _, conflict := range versions {
				cmd.Printf("⚠️  Version conflict detected for package '%s':\n", conflict.Package)

				for version, projects := range conflict.ProjectsWithVersions {
					cmd.Printf("  - Version '%s' used in projects:\n", version)
					for _, project := range projects {
						cmd.Printf("    - %s\n", project)
					}
				}
				cmd.Println()
			}
			if len(versions) == 0 {
				cmd.Println("✅ No dependency version conflicts detected.")
			}
		},
	}

}
