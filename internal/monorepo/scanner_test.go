package monorepo

import (
	"dhakalu/monofile/internal/parsers"
	"testing"
	"testing/fstest"
)

func TestIsProjectFile(t *testing.T) {

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"nested package.json works", "webapp/package.json", true},
		{"dies not match package.json in the middle of a name", "package.json.bak", false},
		{"nested csproj works", "long/long/directory/project.csproj", true},
		{"root pyproject.toml works", "pyproject.toml", true},
		{"does not match name csproj", "csproj", false},
		{"ignores backup files", "project.csproj.bak", false},
		{"nested pyproject.toml is matched", "nested/pyproject.toml", true},
		{"other files are not matched", "random.txt", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isProjectFile(test.path)
			if result != test.expected {
				t.Errorf("For path %s, expected %v but got %v", test.path, test.expected, result)
			}
		})
	}

}

func TestGetProjectConfigurationForPath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		raw       string
		expected  parsers.ProjectConfiguration
		expectErr bool
	}{
		{
			name: "valid csproj file",
			path: "example/project-1/file.csproj",
			raw: `
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
	</Project>
	`,
			expectErr: false,
			expected: parsers.ProjectConfiguration{
				Name: "",
				Path: "example/project-1/file.csproj",
				Language: parsers.Dotnet,
				Dependencies: []parsers.Dependency{
					{
						Name:    "Newtonsoft.Json",
						Version: "13.0.3",
						Type:    parsers.DependencyTypeExternal,
					},
					{
						Name: "example/project-2/project-2.csproj",
						Type: parsers.DependencyTypeInternal,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockFS := fstest.MapFS{
				test.path: &fstest.MapFile{
					Data: []byte(test.raw),
				},
			}

			result, err := getProjectConfigurationForPath(mockFS, test.path)
			if test.expectErr {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Path != test.expected.Path {
				t.Errorf("expected path %s, got %s", test.expected.Path, result.Path)
			}

			if result.Language != test.expected.Language {
				t.Errorf("expected language %s, got %s", test.expected.Language, result.Language)
			}

			if len(result.Dependencies) != len(test.expected.Dependencies) {
				t.Errorf("expected %d dependencies, got %d", len(test.expected.Dependencies), len(result.Dependencies))
			}
		})
	}
}
