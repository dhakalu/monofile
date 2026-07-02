package builders

import "dhakalu/monofile/internal/parsers"

type BuildStatus int

const (
	BuildStatusSuccess BuildStatus = iota
	BuildStatusFailed
)

type BuildInfo struct {
	Status BuildStatus
	Error  error
}

type ProjectBuilder interface {
	// Build the project given the path of the file
	// Project files can be one of:
	// package.json
	// *.csproj
	// *.pytomal
	Build(p parsers.ProjectConfiguration) BuildInfo
}
