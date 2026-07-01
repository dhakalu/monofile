package monorepo

import (
	"dhakalu/monofile/internal/model"
	"dhakalu/monofile/internal/parsers"
	"io/fs"
	"log/slog"
	"path/filepath"
)

// Scanner scans the monorepo and returns a map of projects.

type Scanner struct {
	RepoMetadata model.RepoMetadata
	ProjectsMap  map[string]parsers.ProjectConfiguration
}

func NewScanner(repoMetadata model.RepoMetadata) *Scanner {
	return &Scanner{
		RepoMetadata: repoMetadata,
		ProjectsMap:  make(map[string]parsers.ProjectConfiguration),
	}
}

func (s *Scanner) Scan() error {
	return fs.WalkDir(s.RepoMetadata.FileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			project, err := getLangByFileName(s.RepoMetadata.FileSystem, path)
			if err != nil {
				slog.Error("error parsing project file", "path", path, "error", err)
				return err
			}
			if project != nil {
				s.ProjectsMap[path] = *project
			}
		}
		return nil
	})
}

func getLangByFileName(fs fs.FS, path string) (*parsers.ProjectConfiguration, error) {
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	switch ext {
	case ".csproj":
		parser := parsers.DotnetProjectParser{}
		return parser.Parse(fs, path)
	case ".json":
		if filename == "package.json" {
			return &parsers.ProjectConfiguration{
				Name:     filename,
				Path:     filepath.Dir(filename),
				Language: parsers.Javascript,
			}, nil
		}
		return nil, nil
	default:
		return nil, nil
	}
}
