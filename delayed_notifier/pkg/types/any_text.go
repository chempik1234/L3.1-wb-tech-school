package types

// AnyText is a value type for any text
type AnyText string

// NewAnyText creates a new AnyText from giving text
func NewAnyText(text string) AnyText {
	return AnyText(text)
}

// String returns value of AnyText of type string
func (d AnyText) String() string {
	return string(d)
}
