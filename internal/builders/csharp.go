package builders

import (
	"dhakalu/monofile/internal/parsers"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type DotnetBuilder struct {

}

func (b *DotnetBuilder) Build(p parsers.ProjectConfiguration) BuildInfo {
	slog.Info("Building .NET project", slog.String("name", p.Name), slog.String("path", p.Path))
	cmd := exec.Command("dotnet", "build")
	cmd.Dir = filepath.Dir(p.Path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		slog.Error("Error building .NET project", slog.String("name", p.Name), slog.String("path", p.Path), slog.Any("error", err))
		return BuildInfo{
			Status: BuildStatusFailed,
		}
	}
	return BuildInfo{
		Status: BuildStatusSuccess,
	}
}