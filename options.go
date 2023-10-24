package flat

const (
	// Default delimiter used if none is specified
	DefaultDelimiter string = "."
)

// Options allows customization of input and output behavior
type Options struct {
	// Delimiter used for key concatenation
	Delimiter string
	// Function for modifying keys (optional)
	Fold func(string) string
}

// DefaultOptions creates an option object with default values.
func DefaultOptions(opts ...func(options *Options)) *Options {
	options := &Options{
		Delimiter: DefaultDelimiter,
		Fold:      func(s string) string { return s },
	}

	for _, opt := range opts {
		if opt != nil {
			opt(options)
		}
	}

	return options
}

// WithOptions sets custom options for input and output behavior.
func WithOptions(options *Options) func(o *Options) {
	return func(o *Options) {
		if options != nil {
			*o = *options
		}
	}
}

// WithDelimiter sets a custom delimiter for key concatenation.
func WithDelimiter(delimiter string) func(o *Options) {
	return func(o *Options) {
		o.Delimiter = delimiter
	}
}

// WithFold sets a custom function for modifying keys (optional).
func WithFold(fold func(string) string) func(o *Options) {
	return func(o *Options) {
		o.Fold = fold
	}
}
