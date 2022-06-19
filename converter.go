package markdownconverter

// Converter is an interface that represents the contract used the the various implementations
type Converter interface {
	// Format returns a unique name for the converter
	Format() string
	// Parse will parse the standard markdown and return the converted data
	Parse(markdown []byte) ([]byte, error)
}
