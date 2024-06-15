package route

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/server"
)

type productHandler struct {
	log            *slog.Logger
	productService contract.ProductService
	jwtService     contract.JWTService
}

func NewProductHandler(log *slog.Logger, productService contract.ProductService, jwtService contract.JWTService) *productHandler {
	return &productHandler{
		log:            log,
		productService: productService,
		jwtService:     jwtService,
	}
}

func (h *productHandler) RegisterRoute(s *server.Server) {
	mx := middleware.New(h.jwtService, h.log)

	s.Use(mx.WithAuth)
	s.Group(func(s *server.Server) {
		s.Post("/products", h.create)
		s.Delete("/products/{id}", h.delete)
	})
}

func (h *productHandler) create(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateProduct

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		restutil.NewBadRequestError(w, err.Error())
		return
	}

	// TODO: improve error handling
	err = h.productService.Create(payload)
	if err != nil {
		if errors.Is(err, product.ErrInvalidPrice) {
			restutil.NewUnprocessableEntityError(w, err.Error())
			return
		}
		if errors.Is(err, product.ErrInvalidQuantity) {
			restutil.NewUnprocessableEntityError(w, err.Error())
			return
		}
		if errors.Is(err, product.ErrEmptyProductName) {
			restutil.NewUnprocessableEntityError(w, err.Error())
			return
		}
		if errors.Is(err, service.ErrProductNameAlreadyTaken) {
			restutil.NewConflictError(w, err.Error())
			return
		}
		restutil.NewInternalServerError(w, "cannot create resource")
		return
	}

	restutil.WriteSuccess(w, http.StatusCreated)
}

func (h *productHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedID, _ := strconv.ParseInt(id, 10, 64)

	err := h.productService.Delete(int(parsedID))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			restutil.NewBadRequestError(w, err.Error())
			return
		}
		restutil.NewInternalServerError(w, "cannot delete resource")
		return
	}

	restutil.WriteSuccess(w, http.StatusOK)
}
