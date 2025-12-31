package v1handler

import (
	"net/http"
	v1dto "shopify/internal/dto/v1"
	v1service "shopify/internal/service/v1"
	"shopify/internal/utils"
	"shopify/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service v1service.UserService
}

func NewUserHandler(service v1service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUsers(ctx *gin.Context) {
	var params v1dto.GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	users, total, err := uh.service.GetAllUsers(ctx, params.Search, params.OrderBy, params.Sort, params.Page, params.Limit, params.IsSoftDeleted)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	usersDTO := v1dto.MapUsersToDTO(users)

	paginationResp := utils.NewPaginationResponse(usersDTO, params.Page, params.Limit, total)

	utils.ResponseSuccess(ctx, http.StatusOK, paginationResp)
}

func (uh *UserHandler) GetUserByUUID(ctx *gin.Context) {
	var params v1dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	user, err := uh.service.GetUserByUUID(ctx, userUuid)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(user)

	utils.ResponseSuccess(ctx, http.StatusOK, userDTO)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input v1dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapCreateInputToModel()

	createdUser, err := uh.service.CreateUser(ctx, user)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(createdUser)

	utils.ResponseSuccess(ctx, http.StatusCreated, userDTO)
}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var params v1dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	var input v1dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapUpdateInputToModel(userUuid)

	updatedUser, err := uh.service.UpdateUser(ctx, user)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(updatedUser)

	utils.ResponseSuccess(ctx, http.StatusOK, userDTO)
}

func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)

	err = uh.service.DeleteUser(ctx, userUuid)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

func (uh *UserHandler) SoftDeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	err = uh.service.SoftDeleteUser(ctx, userUuid)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}

func (uh *UserHandler) RestoreUser(ctx *gin.Context) {
	var params v1dto.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userUuid, err := uuid.Parse(params.Uuid)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	restoredUser, err := uh.service.RestoreUser(ctx, userUuid)

	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(restoredUser)

	utils.ResponseSuccess(ctx, http.StatusOK, userDTO)
}
