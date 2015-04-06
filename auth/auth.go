package auth

const (
	validToken   = "3A3E6C4C51F12DF2415682CCF9D18"
	invalidToken = "8A95585DD5B64E33D5BF4C8F4E849"
)

var (
	tokens = map[string]bool{
		validToken:   true,
		invalidToken: false,
	}
)

// Auth is a provider of auth token management/lookup.
// TODO this is stubbed and needs to be DB/cache enabled.
type Auth struct {
	Tokens map[string]bool
}

// New is a factory method that returns an instance of Auth.
func New() *Auth {
	return &Auth{Tokens: tokens}
}

// Validate returns true if the token was found and is valid.
func (t *Auth) Valid(token string) bool {
	a, ok := t.Tokens[token]
	if !ok {
		return false
	}
	return a
}
