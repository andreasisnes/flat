package goflat

const (
	// Default delimiter
	DefaultDelimiter string   = "."
	UnchangedFold    FoldType = iota
	LowerCaseFold
	UpperCaseFold
)

type FoldType int

// Options that customizes the outcome when flatten or unflattens data
type Options struct {
	Delimiter string
	Fold      FoldType
}
