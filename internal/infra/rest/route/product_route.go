package route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/error"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/response"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/server"
)

type productHandler struct {
	log            *slog.Logger
	productService contract.ProductService
	jwtAuth        auth.TokenAuthInterface
}

func NewProductHandler(log *slog.Logger, productService contract.ProductService, jwtAuth auth.TokenAuthInterface) *productHandler {
	return &productHandler{log, productService, jwtAuth}
}

func (h *productHandler) RegisterRoute(s *server.Server) {
	mid := middleware.New(h.jwtAuth, h.log)

	s.Group(func(s *server.Server) {
		s.Use(mid.WithAuth)

		s.Post("/products", h.create)
		s.Get("/products", h.getAll)
		s.Get("/products/{id}", h.getByID)
		s.Get("/products/slug/{slug}", h.getBySlug)
		s.Delete("/products/{id}", h.delete)
	})
}

func (h *productHandler) create(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateProduct

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		error.NewBadRequestError(w, err.Error())
		return
	}

	// TODO: improve error handling
	err = h.productService.Create(payload)
	if err != nil {
		if errors.Is(err, product.ErrInvalidPrice) {
			error.NewUnprocessableEntityError(w, err.Error())
			return
		}
		if errors.Is(err, product.ErrInvalidQuantity) {
			error.NewUnprocessableEntityError(w, err.Error())
			return
		}
		if errors.Is(err, product.ErrEmptyProductName) {
			error.NewUnprocessableEntityError(w, err.Error())
			return
		}
		if errors.Is(err, service.ErrProductNameAlreadyTaken) {
			error.NewConflictError(w, err.Error())
			return
		}
		error.NewInternalServerError(w, "cannot create resource")
		return
	}

	restutil.WriteSuccess(w, http.StatusCreated)
}

func (h *productHandler) delete(w http.ResponseWriter, r *http.Request) {
	publicID := r.PathValue("id")

	err := h.productService.Delete(publicID)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			error.NewBadRequestError(w, err.Error())
			return
		}
		error.NewInternalServerError(w, "cannot delete resource")
		return
	}

	restutil.WriteSuccess(w, http.StatusOK)
}

func (h *productHandler) getBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	p, err := h.productService.GetBySlug(slug)
	if err != nil {
		error.NewBadRequestError(w, err.Error())
		return
	}
	product := response.Product{
		PublicID:  p.PublicID,
		Slug:      p.Slug,
		Name:      p.Name,
		Price:     p.Price,
		Quantity:  p.Quantity,
		Enabled:   p.Enabled,
		CreatedAt: p.CreatedAt,
	}

	restutil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"product": product,
	})
}

func (h *productHandler) getAll(w http.ResponseWriter, r *http.Request) {
	ps, err := h.productService.GetAll()
	if err != nil {
		error.NewBadRequestError(w, err.Error())
		return
	}

	var products []response.Product
	for _, p := range ps {
		products = append(products, response.Product{
			PublicID:  p.PublicID,
			Slug:      p.Slug,
			Name:      p.Name,
			Price:     p.Price,
			Quantity:  p.Quantity,
			Enabled:   p.Enabled,
			CreatedAt: p.CreatedAt,
		})
	}
	restutil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

func (h *productHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	p, err := h.productService.GetByPublicID(id)
	if err != nil {
		error.NewBadRequestError(w, err.Error())
		return
	}
	product := response.Product{
		PublicID:  p.PublicID,
		Slug:      p.Slug,
		Name:      p.Name,
		Price:     p.Price,
		Quantity:  p.Quantity,
		Enabled:   p.Enabled,
		CreatedAt: p.CreatedAt,
	}

	restutil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"product": product,
	})
}
