package route

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/application/service"
	"github.com/bernardinorafael/kn-server/internal/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/response"
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
	mid := middleware.New(h.jwtService, h.log)

	s.Use(mid.WithAuth)
	s.Group(func(s *server.Server) {
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
	publicID := r.PathValue("id")

	err := h.productService.Delete(publicID)
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

func (h *productHandler) getBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	p, err := h.productService.GetBySlug(slug)
	if err != nil {
		restutil.NewBadRequestError(w, err.Error())
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
		restutil.NewBadRequestError(w, err.Error())
		return
	}

	products := []response.Product{}
	for _, p := range ps {
		product := response.Product{
			PublicID:  p.PublicID,
			Slug:      p.Slug,
			Name:      p.Name,
			Price:     p.Price,
			Quantity:  p.Quantity,
			Enabled:   p.Enabled,
			CreatedAt: p.CreatedAt,
		}
		products = append(products, product)
	}

	restutil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

func (h *productHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	p, err := h.productService.GetByPublicID(id)
	if err != nil {
		restutil.NewBadRequestError(w, err.Error())
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
