package handler

import (
	"encoding/json"
	"net/http"

	"github.com/isogram/clean-golang/pkg/bookmark"
	"github.com/isogram/clean-golang/pkg/entity"
	"github.com/isogram/clean-golang/pkg/utils"

	"github.com/labstack/echo"
)

// HTTPBookmarkHandler model
type HTTPBookmarkHandler struct {
	useCase bookmark.UseCase
}

// NewHTTPBookmarkHandler function for initialise *HTTPAuthHandler
func NewHTTPBookmarkHandler(service bookmark.UseCase) *HTTPBookmarkHandler {
	return &HTTPBookmarkHandler{
		useCase: service,
	}
}

func (h *HTTPBookmarkHandler) bookmarkIndex(c echo.Context) error {
	//initialize json schema template pointer
	jsonSchemaTemplate := new(utils.JSONSchemaTemplate)

	errorMessage := "Error reading bookmarks"
	var data []*entity.Bookmark
	var err error
	name := c.QueryParam("name")
	switch {
	case name == "":
		data, err = h.useCase.FindAll()
	default:
		data, err = h.useCase.Search(name)
	}

	if err != nil && err != entity.ErrNotFound {
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

func (h *HTTPBookmarkHandler) bookmarkAdd(c echo.Context) error {
	//initialize json schema template pointer
	jsonSchemaTemplate := new(utils.JSONSchemaTemplate)

	errorMessage := "Error adding bookmark"
	var b *entity.Bookmark
	err := json.NewDecoder(c.Request().Body).Decode(&b)
	if err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = utils.ErrorPayloadInvalid
		jsonSchemaTemplate.Code = http.StatusBadRequest

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	b.ID, err = h.useCase.Store(b)
	if err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = errorMessage
		jsonSchemaTemplate.Code = http.StatusBadRequest

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	jsonSchemaTemplate.Success = true
	jsonSchemaTemplate.Message = "OK"
	jsonSchemaTemplate.Code = http.StatusCreated
	jsonSchemaTemplate.SetData(b)

	return jsonSchemaTemplate.ShowHTTPResponse(c)
}

func (h *HTTPBookmarkHandler) bookmarkFind(c echo.Context) error {
	//initialize json schema template pointer
	jsonSchemaTemplate := new(utils.JSONSchemaTemplate)

	errorMessage := "Error reading bookmark"
	id := c.Param("id")
	data, err := h.useCase.Find(entity.StringToID(id))

	if err != nil && err != entity.ErrNotFound {
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

func (h *HTTPBookmarkHandler) bookmarkDelete(c echo.Context) error {
	//initialize json schema template pointer
	jsonSchemaTemplate := new(utils.JSONSchemaTemplate)

	errorMessage := "Error removing bookmark"
	id := c.Param("id")
	err := h.useCase.Delete(entity.StringToID(id))

	if err != nil {
		jsonSchemaTemplate.Success = false
		jsonSchemaTemplate.Message = errorMessage
		jsonSchemaTemplate.Code = http.StatusBadRequest

		return jsonSchemaTemplate.ShowHTTPResponse(c)
	}

	jsonSchemaTemplate.Success = true
	jsonSchemaTemplate.Message = "OK"
	jsonSchemaTemplate.Code = http.StatusOK
	jsonSchemaTemplate.SetData(utils.Empty{})

	return jsonSchemaTemplate.ShowHTTPResponse(c)
}

//MakeBookmarkHandlers make url handlers
func (h *HTTPBookmarkHandler) MakeBookmarkHandlers(group *echo.Group) {
	group.GET("", h.bookmarkIndex)

	group.POST("", h.bookmarkAdd)

	group.GET("/:id", h.bookmarkFind)

	group.DELETE("/:id", h.bookmarkDelete)
}
