package monorepo

import (
	"dhakalu/monofile/internal/parsers"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

// Scanner scans the monorepo and returns a map of projects.
type Scanner struct {
	FileSystem fs.FS
}

// Option is a function that configures the Scanner.
type Option func(*Scanner)

// WithFS sets the file system to use for scanning.
func WithFS(fileSystem fs.FS) Option {
	return func(s *Scanner) {
		s.FileSystem = fileSystem
	}
}

// NewScanner creates a new Scanner with the given options.
func NewScanner(repoRoot string, opts ...Option) (*Scanner, error) {
	sc := &Scanner{}
	for _, opt := range opts {
		opt(sc)
	}

	if sc.FileSystem == nil {
		if repoRoot == "" {
			cwd, err := os.Getwd()
			if err != nil {
				return nil, err
			}
			repoRoot = cwd
		}
		dfs := os.DirFS(repoRoot)

		sc.FileSystem = dfs
	}
	return sc, nil
}

type ScanResult struct {
	Projects map[string]parsers.ProjectConfiguration
	Warnings []ScanWarning
}

type ScanWarning struct {
	Path    string
	Message string
}

// Scan scans the monorepo and returns a map of projects.
func (s *Scanner) Scan() (*ScanResult, error) {
	result := &ScanResult{
		Projects: make(map[string]parsers.ProjectConfiguration),
	}
	projectsMap := make(map[string]parsers.ProjectConfiguration)
	err := fs.WalkDir(s.FileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && isProjectFile(path) {
			slog.Info("found project file", slog.String("path", path))
			project, perr := getProjectConfigurationForPath(s.FileSystem, path)
			if perr != nil {
				result.Warnings = append(result.Warnings, ScanWarning{
					Path:    path,
					Message: perr.Error(),
				})
				return nil
			}
			if project != nil {
				projectsMap[path] = *project
			} else {
				result.Warnings = append(result.Warnings, ScanWarning{
					Path:    path,
					Message: "prooject configuration could not be parsed",
				})
			}
		}
		return nil
	})
	result.Projects = projectsMap
	return result, err
}

func isProjectFile(path string) bool {
	var exactProjectFiles = map[string]bool{
		"package.json":   true,
		"pyproject.toml": true,
	}
	var extentionBasedProjectFiles = map[string]bool{
		".csproj": true,
	}
	base := filepath.Base(path)
	if exactProjectFiles[base] {
		return true
	}
	ext := filepath.Ext(path)
	return extentionBasedProjectFiles[ext]
}

var ErrNoParserFound = errors.New("no parser found for file")

func getProjectConfigurationForPath(f fs.FS, path string) (*parsers.ProjectConfiguration, error) {
	slog.Debug("getting parser for ", slog.String("project", path))
	parser := parsers.GetParserByFilePath(path)
	if parser == nil {
		return nil, ErrNoParserFound
	}
	return parser.Parse(f, path)
}
