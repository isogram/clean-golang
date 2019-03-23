package handler

import (
	"net/http"

	"github.com/isogram/clean-golang/pkg/entity"
	"github.com/isogram/clean-golang/pkg/middleware"
	"github.com/isogram/clean-golang/pkg/user"
	"github.com/isogram/clean-golang/pkg/utils"

	"github.com/labstack/echo"
)

// HTTPUserHandler model
type HTTPUserHandler struct {
	useCase user.UseCase
}

// NewHTTPUserHandler function for initialise *HTTPAuthHandler
func NewHTTPUserHandler(service user.UseCase) *HTTPUserHandler {
	return &HTTPUserHandler{
		useCase: service,
	}
}

// Me Handler
func (h *HTTPUserHandler) Me(c echo.Context) error {
	//initialize json schema template pointer
	jsonSchemaTemplate := new(utils.JSONSchemaTemplate)

	errorMessage := "Error reading user"
	var (
		data *entity.User
		err  error
	)

	ID := c.Get("UID").(int64)
	data, err = h.useCase.Me(ID)

	if err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = errorMessage
		jsonSchemaTemplate.Code = http.StatusBadRequest

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	if data == nil {
		jsonSchemaTemplate.Success = true
		jsonSchemaTemplate.Message = errorMessage
		jsonSchemaTemplate.Code = http.StatusNotFound
		jsonSchemaTemplate.SetData(utils.Empty{})

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	jsonSchemaTemplate.Success = true
	jsonSchemaTemplate.Message = "OK"
	jsonSchemaTemplate.Code = http.StatusOK
	jsonSchemaTemplate.SetData(data)

	return jsonSchemaTemplate.ShowHTTPResponse(c)
}

// Register handler
func (h *HTTPUserHandler) Register(c echo.Context) error {
	//initialize json schema template pointer
	jsonSchemaTemplate := new(utils.JSONSchemaTemplate)

	u := new(entity.UserRegister)
	errorMessage := "Error registering user"

	var (
		err error
	)

	if err = c.Bind(u); err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = utils.ErrorPayloadInvalid
		jsonSchemaTemplate.Code = http.StatusBadRequest

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	data, err := h.useCase.Register(u)
	if err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = errorMessage
		jsonSchemaTemplate.Code = http.StatusBadRequest

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	jsonSchemaTemplate.Success = true
	jsonSchemaTemplate.Message = "OK"
	jsonSchemaTemplate.Code = http.StatusCreated
	jsonSchemaTemplate.SetData(data)

	return jsonSchemaTemplate.ShowHTTPResponse(c)
}

// Login handler
func (h *HTTPUserHandler) Login(c echo.Context) error {
	//initialize json schema template pointer
	jsonSchemaTemplate := new(utils.JSONSchemaTemplate)
	response := new(entity.UserLoginResponse)
	u := new(entity.UserLogin)

	var (
		err error
	)

	if err = c.Bind(u); err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = utils.ErrorPayloadInvalid
		jsonSchemaTemplate.Code = http.StatusBadRequest

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	response.User, response.Token, response.RefreshToken, err = h.useCase.Login(u.Email, u.Password)
	if err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = err.Error()
		jsonSchemaTemplate.Code = http.StatusUnauthorized
		jsonSchemaTemplate.SetData(utils.Empty{})

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	response.TokenType = "Bearer"

	jsonSchemaTemplate.Success = true
	jsonSchemaTemplate.Message = "OK"
	jsonSchemaTemplate.Code = http.StatusOK
	jsonSchemaTemplate.SetData(response)

	return jsonSchemaTemplate.ShowHTTPResponse(c)
}

// RefreshToken handler
func (h *HTTPUserHandler) RefreshToken(c echo.Context) error {
	//initialize json schema template pointer
	jsonSchemaTemplate := new(utils.JSONSchemaTemplate)
	response := new(entity.UserLoginResponse)
	u := new(entity.UserRefreshToken)

	var (
		err error
	)

	if err = c.Bind(u); err != nil || u.RefreshToken == "" {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = utils.ErrorPayloadInvalid
		jsonSchemaTemplate.Code = http.StatusBadRequest

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	UID := c.Get("UID").(int64)
	response.User, response.Token, response.RefreshToken, err = h.useCase.RefreshToken(UID, u.RefreshToken)
	if err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = err.Error()
		jsonSchemaTemplate.Code = http.StatusUnauthorized

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	response.TokenType = "Bearer"

	jsonSchemaTemplate.Success = true
	jsonSchemaTemplate.Message = "OK"
	jsonSchemaTemplate.Code = http.StatusOK
	jsonSchemaTemplate.SetData(response)

	return jsonSchemaTemplate.ShowHTTPResponse(c)
}

//MakeUserHandlers make url handlers
func (h *HTTPUserHandler) MakeUserHandlers(group *echo.Group) {
	group.GET("/me", h.Me, middleware.BearerVerify(false))
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.POST("/refresh_token", h.RefreshToken, middleware.BearerVerify(true))
}
