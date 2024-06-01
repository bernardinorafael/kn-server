package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/http/httperror"
)

type productHandler struct {
	l              *slog.Logger
	productService contract.ProductService
}

func NewProductHandler(l *slog.Logger, productService contract.ProductService) *productHandler {
	return &productHandler{l: l, productService: productService}
}

func (h *productHandler) RegisterRoute(mux *http.ServeMux) {
	mux.HandleFunc("POST /product", h.create)
}

func (h *productHandler) create(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateProduct

	if r.Body == nil {
		httperror.NewBadRequestError(w, "cannot parse body")
		return
	}
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		httperror.NewInternalServerError(w, "an unknown error occurred")
		return
	}

	err = h.productService.Create(payload)
	if err != nil {
		httperror.NewInternalServerError(w, "cannot create resource")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
