package userroute

import (
	"net/http"
	"sync"

	httperr "github.com/bernardinorafael/kn-server/helper/error"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
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

	u, err := h.accService.GetByID(id)
	if err != nil {
		httperr.NewNotFoundError(c, err.Error())
		return
	}

	user := UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Document:  u.Document,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	c.JSON(http.StatusOK, user)
}

func (h UserHandler) GetManyUsers(c *gin.Context) {
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
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, users)
}
