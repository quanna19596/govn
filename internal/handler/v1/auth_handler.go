package v1handler

import (
	"log"
	"net/http"
	v1dto "shopify/internal/dto/v1"
	v1service "shopify/internal/service/v1"
	"shopify/internal/utils"
	"shopify/internal/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service v1service.AuthService
}

func NewAuthHandler(service v1service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var input v1dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiresIn, err := ah.service.Login(ctx, input.Email, input.Password)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	response := v1dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	log.Println(response)

	utils.ResponseSuccess(ctx, http.StatusOK, response)
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {
	var input v1dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	if err := ah.service.Logout(ctx, input.RefreshToken); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
}

func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	var input v1dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiresIn, err := ah.service.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	response := v1dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Refresh token generate successfully", response)
}
