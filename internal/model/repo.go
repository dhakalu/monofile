package model

import "io/fs"

// Metadata about the current repository
type RepoMetadata struct {
	FileSystem fs.FS
	Root       string `json:"root"`
}
