package gormrepo

import (
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/team"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"gorm.io/gorm"
)

type teamRepo struct {
	db *gorm.DB
}

func NewTeamRepo(db *gorm.DB) contract.TeamRepository {
	return &teamRepo{db}
}

func (tr teamRepo) Create(t team.Team) (gormodel.Team, error) {
	// TODO implement me
	panic("implement me")
}

func (tr teamRepo) Update(t team.Team) (gormodel.Team, error) {
	// TODO implement me
	panic("implement me")
}

func (tr teamRepo) Delete(publicID string) error {
	// TODO implement me
	panic("implement me")
}

func (tr teamRepo) GetByID(publicID string) (gormodel.Team, error) {
	// TODO implement me
	panic("implement me")
}
