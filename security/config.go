package security

const (
	upperRegex   = `[A-Z]`
	lowerRegex   = `[a-z]`
	digitRegex   = `[0-9]`
	specialRegex = `[^a-zA-Z0-9]`
)

type PasswordConfig struct {
	MinimumChars     int
	NeedUppercase    bool
	NeedDigits       bool
	NeedSpecialChars bool
	NeedLowercase    bool
}

type UserConfig struct {
	PasswordConfig  PasswordConfig
	NeedsActivation bool
}
