package main

import (
	"dhakalu/monofile/internal/model"
	mr "dhakalu/monofile/internal/monorepo"
	"log/slog"
	"os"
)

func main() {
	println("Welcome to monofile")
	dfs := os.DirFS("/Users/udhakal/source/monofile")

	
	scanner := mr.NewScanner(model.RepoMetadata{
		FileSystem: dfs,
		Root:       "/Users/udhakal/source/monofile",
	})
	scanner.Scan()
	for k, v := range scanner.ProjectsMap {
		slog.Info("project", slog.String("name", k), slog.AnyValue(v))
	}

}
