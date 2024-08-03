package route

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/application/service"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/product"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/http/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/http/response"
	"github.com/bernardinorafael/kn-server/internal/infra/http/routeutils"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

const (
	maxImageSize = 5 * 1024 * 1024 // 5mb
)

type productHandler struct {
	log            logger.Logger
	productService contract.ProductService
	jwtAuth        auth.TokenAuthInterface
}

func NewProductHandler(log logger.Logger, productService contract.ProductService, jwtAuth auth.TokenAuthInterface) *productHandler {
	return &productHandler{log, productService, jwtAuth}
}

func (h *productHandler) RegisterRoute(r *chi.Mux) {
	m := middleware.NewWithAuth(h.jwtAuth, h.log)

	r.Route("/products", func(r chi.Router) {
		r.With(m.WithAuth)

		r.Post("/", h.create)

		r.Put("/{id}/price", h.updatePrice)
		r.Put("/{id}/status", h.changeStatus)
		r.Put("/{id}/quantity", h.increaseQuantity)

		r.Get("/", h.getAll)
		r.Get("/{id}", h.getByID)
		r.Get("/slug/{slug}", h.getBySlug)

		r.Delete("/{id}", h.delete)

	})
}

func (h *productHandler) changeStatus(w http.ResponseWriter, r *http.Request) {
	var input dto.ChangeStatus

	err := routeutils.ParseBodyRequest(r, &input)
	if err != nil {
		routeutils.NewBadRequestError(w, "error parsing body request")
		return
	}

	err = h.productService.ChangeStatus(r.PathValue("id"), input.Status)
	if err != nil {
		routeutils.NewBadRequestError(w, err.Error())
		return
	}

	routeutils.WriteSuccessResponse(w, http.StatusCreated)
}

func (h *productHandler) increaseQuantity(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateQuantity

	err := routeutils.ParseBodyRequest(r, &input)
	if err != nil {
		routeutils.NewBadRequestError(w, "error parsing body request")
		return
	}

	err = h.productService.IncreaseQuantity(r.PathValue("id"), input.Amount)
	if err != nil {
		routeutils.NewBadRequestError(w, err.Error())
		return
	}

	routeutils.WriteSuccessResponse(w, http.StatusCreated)
}

func (h *productHandler) updatePrice(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdatePrice

	err := routeutils.ParseBodyRequest(r, &input)
	if err != nil {
		routeutils.NewBadRequestError(w, "error parsing body request")
		return
	}

	err = h.productService.UpdatePrice(r.PathValue("id"), input.Amount)
	if err != nil {
		routeutils.NewBadRequestError(w, err.Error())
		return
	}

	routeutils.WriteSuccessResponse(w, http.StatusCreated)
}

func (h *productHandler) create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	quantity := r.FormValue("quantity")
	price := r.FormValue("price")

	f, fh, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			routeutils.NewUnprocessableEntityError(w, "product image cannot be empty")
			return
		}
		routeutils.NewBadRequestError(w, "cannot parse multipart form")
		return
	}
	defer f.Close()

	// verify if image size is larger than 5mb
	if fh.Size > maxImageSize {
		routeutils.NewBadRequestError(w, "image size too long")
		return
	}

	parsedPrice, _ := strconv.ParseFloat(price, 64)
	parsedQuantity, _ := strconv.Atoi(quantity)

	input := dto.CreateProduct{
		Name:      name,
		Price:     parsedPrice,
		Quantity:  parsedQuantity,
		Image:     f,
		ImageName: fh.Filename,
	}

	err = h.productService.Create(input)
	if err != nil {
		switch {
		case errors.Is(err, product.ErrInvalidPrice):
			routeutils.NewUnprocessableEntityError(w, err.Error())
		case errors.Is(err, product.ErrInvalidQuantity):
			routeutils.NewUnprocessableEntityError(w, err.Error())
		case errors.Is(err, product.ErrEmptyProductName):
			routeutils.NewUnprocessableEntityError(w, err.Error())
		case errors.Is(err, service.ErrProductNameAlreadyTaken):
			routeutils.NewConflictError(w, err.Error())
		default:
			routeutils.NewInternalServerError(w, "cannot create resource")
		}
		return
	}

	routeutils.WriteSuccessResponse(w, http.StatusOK)
}

func (h *productHandler) delete(w http.ResponseWriter, r *http.Request) {
	err := h.productService.Delete(r.PathValue("id"))
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			routeutils.NewBadRequestError(w, err.Error())
			return
		}
		routeutils.NewInternalServerError(w, "cannot delete resource")
		return
	}

	routeutils.WriteSuccessResponse(w, http.StatusOK)
}

func (h *productHandler) getBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	p, err := h.productService.GetBySlug(slug)
	if err != nil {
		routeutils.NewBadRequestError(w, err.Error())
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

	routeutils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"product": product,
	})
}

func (h *productHandler) getAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// TODO: make a parser method to query params
	disabled, err := strconv.ParseBool(query.Get("disabled"))
	if !query.Has("disabled") {
		routeutils.NewBadRequestError(w, "missing [disabled] query parameter")
		return
	}
	if err != nil {
		routeutils.NewBadRequestError(w, "cannot parse [disabled] query params")
		return
	}

	orderBy := query.Get("order_by")
	if !query.Has("order_by") {
		routeutils.NewBadRequestError(w, "missing [order_by] query parameter")
		return
	}

	allProducts, err := h.productService.GetAll(disabled, orderBy)
	if err != nil {
		routeutils.NewBadRequestError(w, err.Error())
		return
	}
	
	var products []response.Product
	for _, p := range allProducts {
		products = append(products, response.Product{
			PublicID:  p.PublicID,
			Slug:      p.Slug,
			Name:      p.Name,
			Price:     p.Price,
			Image:     p.Image,
			Quantity:  p.Quantity,
			Enabled:   p.Enabled,
			CreatedAt: p.CreatedAt,
		})
	}

	routeutils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

func (h *productHandler) getByID(w http.ResponseWriter, r *http.Request) {
	p, err := h.productService.GetByPublicID(r.PathValue("id"))
	if err != nil {
		routeutils.NewBadRequestError(w, err.Error())
		return
	}

	product := response.Product{
		PublicID:  p.PublicID,
		Slug:      p.Slug,
		Name:      p.Name,
		Image:     p.Image,
		Price:     p.Price,
		Quantity:  p.Quantity,
		Enabled:   p.Enabled,
		CreatedAt: p.CreatedAt,
	}

	routeutils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"product": product,
	})
}
