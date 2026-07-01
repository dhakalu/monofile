package parsers

import "path/filepath"

func GetParserByFilePath(path string) Parser {
	ext := filepath.Ext(path)
	switch ext {
	case ".csproj":
		return &DotnetProjectParser{}
	// case ".json":
	// 	if isPackageJson(path) {
	// 		return &JavascriptProjectParser{}, nil
	// 	}
	// 	return nil, nil
	default:
		return nil
	}
}

func isPackageJson(path string) bool {
	return filepath.Base(path) == "package.json"
}