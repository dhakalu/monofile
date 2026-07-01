package parsers

import (
	"encoding/xml"
	"io/fs"
	"path/filepath"
)

// Parses the .csproj to support dotnet operations
type Project struct {
	XMLName        xml.Name        `xml:"Project"`
	Sdk            string          `xml:"Sdk,attr"`
	RootNamespace  string          `xml:"RootNamespace"`
	PropertyGroups []PropertyGroup `xml:"PropertyGroup"`
	ItemGroups     []ItemGroup     `xml:"ItemGroup"`
}

type PropertyGroup struct {
	OutputType      string `xml:"OutputType,omitempty"`
	TargetFramework string `xml:"TargetFramework,omitempty"`
	ImplicitUsings  string `xml:"ImplicitUsings,omitempty"`
	Nullable        string `xml:"Nullable,omitempty"`
}

type ItemGroup struct {
	PackageReferences []PackageReference `xml:"PackageReference"`
	ProjectReferences []ProjectReference `xml:"ProjectReference"`
}

type PackageReference struct {
	Include string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

type ProjectReference struct {
	Include string `xml:"Include,attr"`
}

type DotnetProjectParser struct {
}

func (dps DotnetProjectParser) Parse(fileSystem fs.FS, path string) (*ProjectConfiguration, error) {
	content, err := fs.ReadFile(fileSystem, path)
	if err != nil {
		return nil, err
	}
	var project Project
	err = xml.Unmarshal(content, &project)
	if err != nil {
		return nil, err
	}
	deps := make([]Dependency, 0)
	for _, ig := range project.ItemGroups {

		for _, pr := range ig.ProjectReferences {
			deps = append(deps, Dependency{
				Name: filepath.Join(filepath.Dir(path), pr.Include),
				Type: DependencyTypeInternal,
			})
		}

		for _, pkg := range ig.PackageReferences {
			deps = append(deps, Dependency{
				Name:    pkg.Include,
				Type:    DependencyTypeExternal,
				Version: pkg.Version,
			})
		}

	}

	// Convert the parsed project to a ProjectConfiguration
	config := &ProjectConfiguration{
		Name:         project.RootNamespace,
		Path:         path,
		Dependencies: deps,
		Language:     DotNet,
	}

	return config, nil
}
