package monorepo

import (
	"dhakalu/monofile/internal/parsers"
	"maps"
	"slices"
	"strings"
)

type VersionConflictReport struct {
	Language             parsers.SourceLang
	Package              string
	ProjectsWithVersions map[string][]string
}

type VersionMap map[string]map[string][]string

func GetDependencyVersions(projects map[string]parsers.ProjectConfiguration) []VersionConflictReport {
	versionsMap := make(VersionMap)
	for _, project := range projects {
		for _, dep := range project.Dependencies {
			key := strings.Join([]string{string(project.Language), dep.Name}, ":")
			if _, exists := versionsMap[key]; !exists {
				versionsMap[key] = make(map[string][]string)
			}

			if _, exists := versionsMap[key][dep.Version]; !exists {
				versionsMap[key][dep.Version] = []string{}
			}
			versionsMap[key][dep.Version] = append(versionsMap[key][dep.Version], project.Path)
		}
	}
	var conflicts []VersionConflictReport
	for key, v := range versionsMap {
		versions := slices.Collect(maps.Keys(versionsMap[key]))
		if len(versions) > 1 {
			parts := strings.Split(key, ":")
			conflicts = append(conflicts, VersionConflictReport{
				Language:             parsers.SourceLang(parts[0]),
				Package:              parts[1],
				ProjectsWithVersions: v,
			})
		}
	}
	return conflicts
}
