package usersHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jetsadawwts/go-restapi/config"
	"github.com/jetsadawwts/go-restapi/modules/entities"
	"github.com/jetsadawwts/go-restapi/modules/users"
	"github.com/jetsadawwts/go-restapi/modules/users/usersUsecases"
	"github.com/jetsadawwts/go-restapi/pkg/auth"
)

type userHandlerErrCode string

const (
	SignUpCustomerErr     userHandlerErrCode = "users-001"
	SignInErr             userHandlerErrCode = "users-002"
	RefreshPassportErr    userHandlerErrCode = "users-003"
	SignOutErr            userHandlerErrCode = "users-004"
	SignUpAdminErr        userHandlerErrCode = "users-005"
	GenerateAdminTokenErr userHandlerErrCode = "users-006"
	GetUserProfileErr     userHandlerErrCode = "users-007"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	RefreshPassport(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	GenerateAdminToken(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg          config.IConfig
	usersUsecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usersUsecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:          cfg,
		usersUsecase: usersUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	//Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignUpCustomerErr),
			err.Error(),
		).Res()
	}

	//Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	//Insert
	result, err := h.usersUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpCustomerErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
	//Request body parser
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignUpCustomerErr),
			err.Error(),
		).Res()
	}

	//Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	//Insert
	result, err := h.usersUsecase.InsertAdmin(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpAdminErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpAdminErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(SignUpAdminErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) GenerateAdminToken(c *fiber.Ctx) error {
	adminToken, err := auth.NewAuth(auth.Admin, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code,
			string(GenerateAdminTokenErr), err.Error()).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, &struct {
		Token string `json: "token"`
	}{
		Token: adminToken.SignToken(),
	}).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignInErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignInErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) RefreshPassport(c *fiber.Ctx) error {
	req := new(users.UserRefreshCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RefreshPassportErr),
			err.Error(),
		).Res()
	}

	passport, err := h.usersUsecase.RefreshPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RefreshPassportErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, passport).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignOutErr),
			err.Error(),
		).Res()
	}

	if err := h.usersUsecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(SignOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	result, err := h.usersUsecase.GetUserProfile(userId)

	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(GetUserProfileErr),
				err.Error()).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(GetUserProfileErr),
				err.Error()).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
