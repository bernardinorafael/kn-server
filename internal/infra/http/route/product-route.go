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
	"github.com/bernardinorafael/kn-server/internal/infra/http/error"
	"github.com/bernardinorafael/kn-server/internal/infra/http/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/http/response"
	"github.com/bernardinorafael/kn-server/internal/infra/http/restutil"
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
	publicID := r.PathValue("id")

	err := restutil.ParseBody(r, &input)
	if err != nil {
		error.NewBadRequestError(w, "error parsing body request")
		return
	}

	err = h.productService.ChangeStatus(publicID, input.Status)
	if err != nil {
		error.NewBadRequestError(w, err.Error())
		return
	}
	restutil.WriteSuccess(w, http.StatusCreated)
}

func (h *productHandler) increaseQuantity(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdateQuantity
	publicID := r.PathValue("id")

	err := restutil.ParseBody(r, &input)
	if err != nil {
		error.NewBadRequestError(w, "error parsing body request")
		return
	}

	err = h.productService.IncreaseQuantity(publicID, input.Amount)
	if err != nil {
		error.NewBadRequestError(w, err.Error())
		return
	}

	restutil.WriteSuccess(w, http.StatusCreated)
}

func (h *productHandler) updatePrice(w http.ResponseWriter, r *http.Request) {
	var input dto.UpdatePrice
	publicID := r.PathValue("id")

	err := restutil.ParseBody(r, &input)
	if err != nil {
		error.NewBadRequestError(w, "error parsing body request")
		return
	}

	err = h.productService.UpdatePrice(publicID, input.Amount)
	if err != nil {
		error.NewBadRequestError(w, err.Error())
		return
	}

	restutil.WriteSuccess(w, http.StatusCreated)
}

func (h *productHandler) create(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	quantity := r.FormValue("quantity")
	price := r.FormValue("price")

	f, fh, err := r.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			error.NewUnprocessableEntityError(w, "product image cannot be empty")
			return
		}
		error.NewBadRequestError(w, "cannot parse multipart form")
		return
	}
	defer f.Close()

	// verify if image size is larger than 5mb
	if fh.Size > maxImageSize {
		error.NewBadRequestError(w, "image size too long")
		return
	}

	parsedPrice, _ := strconv.ParseFloat(price, 64)
	parsedQuantity, _ := strconv.Atoi(quantity)

	input := dto.CreateProduct{
		Name:      name,
		Price:     parsedPrice,
		Quantity:  int32(parsedQuantity),
		Image:     f,
		ImageName: fh.Filename,
	}

	err = h.productService.Create(input)
	if err != nil {
		switch {
		case errors.Is(err, product.ErrInvalidPrice):
			error.NewUnprocessableEntityError(w, err.Error())
		case errors.Is(err, product.ErrInvalidQuantity):
			error.NewUnprocessableEntityError(w, err.Error())
		case errors.Is(err, product.ErrEmptyProductName):
			error.NewUnprocessableEntityError(w, err.Error())
		case errors.Is(err, service.ErrProductNameAlreadyTaken):
			error.NewConflictError(w, err.Error())
		default:
			error.NewInternalServerError(w, "cannot create resource")
		}
		return
	}
	restutil.WriteSuccess(w, http.StatusOK)
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
	query := r.URL.Query()

	// TODO: make a parser method to query params
	disabled, err := strconv.ParseBool(query.Get("disabled"))
	if !query.Has("disabled") {
		error.NewBadRequestError(w, "missing [disabled] query parameter")
		return
	}
	if err != nil {
		error.NewBadRequestError(w, "cannot parse [disabled] query params")
		return
	}

	orderBy := query.Get("order_by")
	if !query.Has("order_by") {
		error.NewBadRequestError(w, "missing [order_by] query parameter")
		return
	}

	ps, err := h.productService.GetAll(disabled, orderBy)
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
			Image:     p.Image,
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
		Image:     p.Image,
		Price:     p.Price,
		Quantity:  p.Quantity,
		Enabled:   p.Enabled,
		CreatedAt: p.CreatedAt,
	}

	restutil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"product": product,
	})
}