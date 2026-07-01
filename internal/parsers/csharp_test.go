package parsers

import (
	"testing"
	"testing/fstest"
)

func TestGetProjectFromPath(t *testing.T) {

	xmlData := `
	<Project Sdk="Microsoft.NET.Sdk">
	  <PropertyGroup>
	    <OutputType>Exe</OutputType>
	    <TargetFramework>net9.0</TargetFramework>
	  </PropertyGroup>
	  <ItemGroup>
	    <PackageReference Include="Newtonsoft.Json" Version="13.0.3" />
	  </ItemGroup>
	  <ItemGroup>
	    <ProjectReference Include="../project-2/project-2.csproj" />
	  </ItemGroup>
	</Project>`

	mockFS := fstest.MapFS{
		"example/project-1/file.csproj": &fstest.MapFile{
			Data: []byte(xmlData),
		},
	}

	expected := ProjectConfiguration{
		Name:     "",
		Path:     "example/project-1/file.csproj",
		Language: Dotnet,
		Dependencies: []Dependency{
			{
				Name:    "Newtonsoft.Json",
				Type:    DependencyTypeExternal,
				Version: "13.0.3",
			},
			{
				Name: "example/project-2/project-2.csproj",
				Type: DependencyTypeInternal,
			},
		},
	}

	dps := DotnetProjectParser{}
	data, err := dps.Parse(mockFS, "example/project-1/file.csproj")

	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if data.Name != expected.Name {
		t.Errorf("expected name %s, got %s", expected.Name, data.Name)
	}

	if data.Path != expected.Path {
		t.Errorf("expected path %s, got %s", expected.Path, data.Path)
	}

	if data.Language != expected.Language {
		t.Errorf("expected language %v, got %v", expected.Language, data.Language)
	}

	if len(data.Dependencies) != len(expected.Dependencies) {
		t.Fatalf("expected %d dependencies, got %d", len(expected.Dependencies), len(data.Dependencies))
	}

	for i, dep := range data.Dependencies {
		if dep.Name != expected.Dependencies[i].Name {
			t.Errorf("expected dependency name %s, got %s", expected.Dependencies[i].Name, dep.Name)
		}
		if dep.Type != expected.Dependencies[i].Type {
			t.Errorf("expected dependency type %v, got %v", expected.Dependencies[i].Type, dep.Type)
		}
	}

}
