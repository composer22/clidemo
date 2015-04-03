package auth

var (
	tokens = map[string]bool{
		"3A3E6C4C51F12DF2415682CCF9D18": true,
		"8A95585DD5B64E33D5BF4C8F4E849": false,
	}
)

// Auth is a provider of auth token management/lookup.
// TODO Right now this is stubbed, but will eventally need to be DB/cache enabled.
type Auth struct {
	Tokens map[string]bool
}

// New is a factory method that returns an instance of Auth
func New() *Auth {
	return &Auth{Tokens: tokens}
}

// Validate returns true if the token was found to be valid
func (t *Auth) Valid(token string) bool {
	a, ok := t.Tokens[token]
	if !ok {
		return false
	}
	return a
}
