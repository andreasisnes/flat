package goflat

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
