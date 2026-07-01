package parsers

import "io/fs"

// Parsers are responsible for parsing project files in a monorepo
// and extracting relevant information about the projects, their
// dependencies, and configurations. Each parser is designed to
// handle specific types of project files (e.g., .csproj for .NET
// projects, package.json for Node.js projects, etc.) and convert
// them into a standardized format that can be used by the monorepo
// scanner.

type SourceLang int

const (
	None       SourceLang = iota
	Python     SourceLang = 1
	DotNet     SourceLang = 2
	Javascript SourceLang = 3
	Go         SourceLang = 4
)
// Configuration of each project in monorepo, including its dependencies and their versions.
// Configuration is parsed from the project files (e.g., .csproj for .NET projects)
// and stored in a map for easy access.
type ProjectConfiguration struct {
	// Name of the project
	Name         string
	Path         string
	Language     SourceLang
	Dependencies []Dependency
}

const (
	DependencyTypeInternal = "internal"
	DependencyTypeExternal = "external"
)

type DependencyType string

type Dependency struct {
	// Name of the package or project that is a dependency
	// For example, in a .NET project, this could be the name of a
	// NuGet package or another project within the monorepo.
	// For internal depnendencies, this could be the path to
	// the project file (e.g., .csproj) of the dependent project.
	Name    string
	Type    DependencyType
	Version string
}

type Parser interface {
	// Parse the project given the path of the file
	// Project files can be one of:
	// package.json
	// *.csproj
	// *.pytomal
	Parse(fileSystem fs.FS, path string) (*ProjectConfiguration, error)
}
