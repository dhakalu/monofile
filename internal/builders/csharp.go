package builders

import (
	"dhakalu/monofile/internal/parsers"
	"os"
	"os/exec"
	"path/filepath"
)

type DotnetBuilder struct{}

func (DotnetBuilder) Build(p parsers.ProjectConfiguration) BuildInfo {
	cmd := exec.Command("dotnet", "build")
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
