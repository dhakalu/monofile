package builders

import (
	"dhakalu/monofile/internal/parsers"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

type DotnetBuilder struct{}

// Build c# projects
func (DotnetBuilder) Build(p parsers.ProjectConfiguration, args map[string]string) BuildInfo {
	cmd := exec.Command("dotnet", "build")
	for arg, value := range args {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--%s", arg), value)
	}
	slog.Info("running", slog.String("command", cmd.String()))
	cmd.Dir = filepath.Dir(p.Path)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return BuildInfo{
			Status: BuildStatusFailed,
			Error:  err,
		}
	}
	return BuildInfo{
		Status: BuildStatusSuccess,
	}
}
