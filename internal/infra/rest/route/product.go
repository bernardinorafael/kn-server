package route

import (
	"log/slog"
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/restutil"
)

type productHandler struct {
	l              *slog.Logger
	productService contract.ProductService
}

func NewProductHandler(l *slog.Logger, productService contract.ProductService) *productHandler {
	return &productHandler{l: l, productService: productService}
}

func (h *productHandler) RegisterRoute(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", h.create)
}

func (h *productHandler) create(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateProduct

	err := restutil.ParseBody(r, &payload)
	if err != nil {
		restutil.NewBadRequestError(w, "cannot parse body")
		return
	}

	err = h.productService.Create(payload)
	if err != nil {
		restutil.NewInternalServerError(w, "cannot create resource")
	}

	restutil.WriteSuccess(w, http.StatusCreated, map[string]interface{}{
		"message": "success",
	})
}
