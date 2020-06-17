package jwt

type JWTError struct {
	Violates map[string][]string
}

func (JWTError) Error() string {
	return "Token is invalid or expired"
}

func (e JWTError) Violations() map[string][]string {
	return e.Violates
}

func (JWTError) Validation() bool {
	return true
}

func (JWTError) ServiceError() bool {
	return false
}
