package v1service

import (
	"shopify/internal/repository"
	"shopify/internal/utils"
	"shopify/pkg/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService) AuthService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
	}
}

func (as *authService) Login(ctx *gin.Context, email string, password string) (string, int, error) {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)
	user, err := as.userRepo.GetByEmail(context, email)

	if err != nil {
		return "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", 0, utils.WrapError(err, "Unable to generate access token", utils.ErrCodeInternalServer)
	}

	return accessToken, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}
