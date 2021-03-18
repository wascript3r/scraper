package auth

type Usecase interface {
	ValidateToken(inputToken string) bool
}
