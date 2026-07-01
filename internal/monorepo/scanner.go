package monorepo

import (
	"dhakalu/monofile/internal/model"
	"dhakalu/monofile/internal/parsers"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// Scanner scans the monorepo and returns a map of projects.

type Scanner struct {
	RepoMetadata model.RepoMetadata
}

func NewScanner(repoRoot string) (*Scanner, error) {
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		repoRoot = cwd
	}
	dfs := os.DirFS(repoRoot)
	md := model.RepoMetadata{
		FileSystem: dfs,
		Root:       repoRoot,
	}
	return &Scanner{
		RepoMetadata: md,
	}, nil
}

func (s *Scanner) Scan() (map[string]parsers.ProjectConfiguration, error) {

	projectsMap := make(map[string]parsers.ProjectConfiguration)
	return projectsMap, fs.WalkDir(s.RepoMetadata.FileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && isProjectFile(path) {
			project, err := getProjectConfigurationForPath(s.RepoMetadata.FileSystem, path)
			if err != nil {
				slog.Error("error parsing project file", "path", path, "error", err)
				return err
			}
			if project != nil {
				projectsMap[path] = *project
			}
		}
		return nil
	})
}

func isProjectFile(path string) bool {
	var projectFiles = []string{
		"package.json",
		"*.csproj",
		"pyproject.toml",
	}
	for _, pattern := range projectFiles {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			slog.Error("error matching pattern", "pattern", pattern, "path", path, "error", err)
			return false
		}
		if matched {
			return true
		}
	}
	return false
}

func getProjectConfigurationForPath(fs fs.FS, path string) (*parsers.ProjectConfiguration, error) {
	parser := parsers.GetParserByFilePath(path)
	if parser == nil {
		slog.Warn("no parser found for file", "path", path)
		return nil, nil
	}
	return parser.Parse(fs, path)
}
