package userroute

import (
	"net/http"
	"sync"

	httperr "github.com/bernardinorafael/kn-server/helper/error"
	"github.com/bernardinorafael/kn-server/helper/validator"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/gin-gonic/gin"
)

var handler *UserHandler
var once sync.Once

type UserHandler struct {
	accService contract.UserService
}

func NewUserHandler(accService contract.UserService) *UserHandler {
	once.Do(func() {
		handler = &UserHandler{accService}
	})
	return handler
}

func (h UserHandler) GetUser(c *gin.Context) {

	id := c.Param("id")
	acc, err := h.accService.GetByID(id)
	if err != nil {
		httperr.NewNotFoundError(c, err.Error())
		return
	}

	account := UserResponse{
		ID:        acc.ID,
		Name:      acc.Name,
		Email:     acc.Email,
		Document:  acc.Document,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}

	c.JSON(http.StatusOK, account)
}

func (h UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	req := &dto.UpdateAccount{}
	if err := c.ShouldBind(req); err != nil {
		httperr.NewBadRequestError(c, "not found/invalid body")
		return
	}

	validations := validator.Validate(req)
	if validations != nil {
		httperr.NewFieldsErrorValidation(c, "invalid fields", validations)
		return
	}

	err := h.accService.UpdateUser(*req, id)
	if err != nil {
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := h.accService.DeleteUser(id)
	if err != nil {
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h UserHandler) GetAccounts(c *gin.Context) {
	allUsers, err := h.accService.GetAll()
	if err != nil {
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	var users []UserResponse
	for _, u := range *allUsers {
		users = append(users, UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Document:  u.Document,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, users)
}
