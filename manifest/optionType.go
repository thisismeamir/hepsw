package manifest

// OptionType defines the type of an option
type OptionType string

const (
	OptionBool   OptionType = "bool"
	OptionString OptionType = "string"
	OptionInt    OptionType = "int"
	OptionEnum   OptionType = "enum"
	OptionList   OptionType = "list"
)
