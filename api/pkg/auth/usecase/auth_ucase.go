package usecase

type AuthUsecase struct {
	token string
}

func New(token string) *AuthUsecase {
	return &AuthUsecase{token}
}

func (a *AuthUsecase) ValidateToken(inputToken string) bool {
	return inputToken == a.token
}
