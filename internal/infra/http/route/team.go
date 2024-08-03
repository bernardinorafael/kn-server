package route

import (
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/http/routeutils"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type teamHandler struct {
	log         logger.Logger
	teamService contract.TeamService
}

func NewTeamHandler(log logger.Logger, teamService contract.TeamService) teamHandler {
	return teamHandler{log, teamService}
}

func (h teamHandler) RegisterRoute(r *chi.Mux) {
	r.Route("/teams", func(r chi.Router) {
		r.Post("/", h.create)
	})
}

func (h teamHandler) create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateTeam

	err := routeutils.ParseBodyRequest(r, &input)
	if err != nil {
		routeutils.NewBadRequestError(w, err.Error())
		return
	}

	err = h.teamService.Create(input)
	if err != nil {
		routeutils.NewBadRequestError(w, "cannot create team")
		return
	}

	routeutils.WriteSuccessResponse(w, http.StatusCreated)
}
