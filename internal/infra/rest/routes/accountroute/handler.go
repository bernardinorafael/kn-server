package accountroute

import (
	"errors"
	"net/http"
	"sync"

	"github.com/bernardinorafael/kn-server/helper/error"
	"github.com/bernardinorafael/kn-server/helper/validator"
	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/gin-gonic/gin"
)

var handler *AccountHandler
var once sync.Once

type AccountHandler struct {
	accService contract.AccountService
}

func NewAccountHandler(accService contract.AccountService) *AccountHandler {
	once.Do(func() {
		handler = &AccountHandler{accService}
	})
	return handler
}

func (h AccountHandler) CreateAccount(c *gin.Context) {
	ctx := restutil.GetContext(c)

	req := dto.CreateAccount{}
	err := c.ShouldBind(&req)
	if err != nil {
		httperr.NewBadRequestError(c, "not found/invalid body")
		return
	}
	validations := validator.Validate(req)
	if validations != nil {
		httperr.NewFieldsErrorValidation(c, "invalid fields", validations)
		return
	}

	err = h.accService.CreateAccount(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			httperr.NewConflictError(c, err.Error())
			return
		}
		httperr.NewConflictError(c, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h AccountHandler) GetUser(c *gin.Context) {
	ctx := restutil.GetContext(c)

	id := c.Param("id")
	acc, err := h.accService.GetByID(ctx, id)
	if err != nil {
		httperr.NewNotFoundError(c, err.Error())
		return
	}

	account := AccountResponse{
		ID:        acc.ID,
		Name:      acc.Name,
		Email:     acc.Email,
		Document:  acc.Document,
		CreatedAt: acc.CreatedAt,
		UpdatedAt: acc.UpdatedAt,
	}

	c.JSON(http.StatusOK, account)
}

func (h AccountHandler) UpdateAccount(c *gin.Context) {
	ctx := restutil.GetContext(c)

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

	err := h.accService.UpdateAccount(ctx, *req, id)
	if err != nil {
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h AccountHandler) DeleteAccount(c *gin.Context) {
	ctx := restutil.GetContext(c)
	id := c.Param("id")

	err := h.accService.DeleteAccount(ctx, id)
	if err != nil {
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h AccountHandler) GetAccounts(c *gin.Context) {
	ctx := restutil.GetContext(c)

	allAcc, err := h.accService.GetAll(ctx)
	if err != nil {
		httperr.NewBadRequestError(c, err.Error())
		return
	}

	var accounts []AccountResponse
	for _, a := range *allAcc {
		accounts = append(accounts, AccountResponse{
			ID:        a.ID,
			Name:      a.Name,
			Email:     a.Email,
			Document:  a.Document,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, accounts)
}
