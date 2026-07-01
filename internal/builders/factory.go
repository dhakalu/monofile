package builders

import "dhakalu/monofile/internal/parsers"

func GetBuilderForLanguage(language parsers.SourceLang) ProjectBuilder {
	switch language {
	case parsers.Dotnet:
		return &DotnetBuilder{}
	default:
		return nil
	}
}
