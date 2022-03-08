package goflat

const (
	// Default delimiter
	DefaultDelimiter = "."
)

// Options that customizes the outcome when flatten or unflattens data
type Options struct {
	Delimiter string
}
