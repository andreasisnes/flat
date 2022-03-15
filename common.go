package goflat

import (
	"fmt"
	"strings"
)

func createDefultOptionsIfNil(options *Options) *Options {
	if options == nil {
		options = &Options{
			Delimiter: DefaultDelimiter,
		}
	}

	if options.Delimiter == "" {
		options.Delimiter = DefaultDelimiter
	}

	return options
}

func getFold(s string, options *Options) string {
	switch options.Fold {
	case UpperCaseFold:
		return strings.ToUpper(s)
	case LowerCaseFold:
		return strings.ToLower(s)
	}

	return s
}

func concatKey(key, idx string, options *Options) string {
	if key == "" {
		return idx
	}

	return fmt.Sprintf("%s%s%s", getFold(key, options), options.Delimiter, getFold(idx, options))
}
