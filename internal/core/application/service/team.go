package service

import (
	"errors"

	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/application/dto"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/team"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"github.com/bernardinorafael/kn-server/pkg/logger"
)

var (
	ErrNotFoundTeam = errors.New("team not found")
)

type teamService struct {
	log      logger.Logger
	teamRepo contract.TeamRepository
}

func NewTeamService(log logger.Logger, teamRepo contract.TeamRepository) contract.TeamService {
	return &teamService{log, teamRepo}
}

func (svc teamService) GetByID(publicID string) (gormodel.Team, error) {
	var team gormodel.Team

	t, err := svc.teamRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("team not found", "public_id", publicID)
		return team, ErrNotFoundTeam
	}

	return t, nil
}

func (svc teamService) Create(dto dto.CreateTeam) error {
	t, err := team.New(dto.OwnerID, dto.Name)
	if err != nil {
		svc.log.Error("error initializing team", "error", err.Error())
		return err
	}

	_, err = svc.teamRepo.Create(*t)
	if err != nil {
		svc.log.Error("error creating team", "error", err.Error())
		return errors.New("cannot create team resource")
	}

	return nil
}
