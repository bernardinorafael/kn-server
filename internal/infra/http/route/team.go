package route

import (
	"net/http"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/http/middleware"
	"github.com/bernardinorafael/kn-server/internal/infra/http/response"
	. "github.com/bernardinorafael/kn-server/internal/infra/http/routeutils"
	"github.com/bernardinorafael/kn-server/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type teamHandler struct {
	log         logger.Logger
	teamService contract.TeamService
	jwtAuth     auth.TokenAuthInterface
}

func NewTeamHandler(log logger.Logger, teamService contract.TeamService, jwtAuth auth.TokenAuthInterface) teamHandler {
	return teamHandler{log, teamService, jwtAuth}
}

func (h teamHandler) RegisterRoute(r *chi.Mux) {
	m := middleware.NewWithAuth(h.jwtAuth, h.log)

	r.Route("/teams", func(r chi.Router) {
		r.Use(m.WithAuth)

		r.Post("/", h.create)
		r.Get("/{id}", h.getByID)
	})
}

func (h teamHandler) getByID(w http.ResponseWriter, r *http.Request) {
	t, err := h.teamService.GetByID(r.PathValue("id"))
	if err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	var members []response.User
	for _, m := range t.Members {
		members = append(members, response.User{
			PublicID:  m.PublicID,
			Name:      m.Name,
			Email:     m.Email,
			Phone:     m.Phone,
			Status:    m.Status,
			CreatedAt: m.CreatedAt,
		})
	}

	team := response.Team{
		PublicID:  t.PublicID,
		Name:      t.Name,
		OwnerID:   t.OwnerID,
		Members:   members,
		CreatedAt: t.CreatedAt,
	}

	WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"team": team,
	})
}

func (h teamHandler) create(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateTeam

	if err := ParseBodyRequest(r, &body); err != nil {
		NewBadRequestError(w, err.Error())
		return
	}

	err := h.teamService.Create(body)
	if err != nil {
		NewBadRequestError(w, "cannot create team")
		return
	}

	WriteSuccessResponse(w, http.StatusCreated)
}
