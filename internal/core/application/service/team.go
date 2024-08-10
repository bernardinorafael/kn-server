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
	ErrTeamNotFound = errors.New("team not found")
)

type teamService struct {
	log      logger.Logger
	teamRepo contract.GormTeamRepository
}

func NewTeamService(log logger.Logger, teamRepo contract.GormTeamRepository) contract.TeamService {
	return &teamService{log, teamRepo}
}

func (svc teamService) GetByID(publicID string) (gormodel.Team, error) {
	team, err := svc.teamRepo.GetByPublicID(publicID)
	if err != nil {
		svc.log.Error("team not found", "public_id", publicID)
		return team, ErrTeamNotFound
	}

	return team, nil
}

func (svc teamService) Create(dto dto.CreateTeam) error {
	t, err := team.New(dto.OwnerID, dto.Name)
	if err != nil {
		svc.log.Error("error initializing team", "error", err.Error())
		return err
	}

	teamModel := gormodel.Team{
		PublicID: t.PublicID(),
		OwnerID:  t.OwnerID(),
		Name:     t.Name(),
	}

	if err = svc.teamRepo.Create(teamModel); err != nil {
		svc.log.Error("error creating team", "error", err.Error())
		return errors.New("cannot create team resource")
	}

	return nil
}
