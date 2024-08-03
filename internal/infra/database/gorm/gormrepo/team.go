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

/*
* TODO: remove entity/model mapping logic from repositories and do it into service layer
 */

func NewTeamRepo(db *gorm.DB) contract.TeamRepository {
	return &teamRepo{db}
}

func (r teamRepo) Create(t team.Team) (gormodel.Team, error) {
	var team gormodel.Team

	newTeam := gormodel.Team{
		PublicID: t.PublicID(),
		OwnerID:  t.OwnerID(),
		Name:     t.Name(),
		Members:  nil,
	}

	err := r.db.
		Create(newTeam).
		First(&team).
		Error
	if err != nil {
		return team, nil
	}

	return team, nil
}

func (r teamRepo) Update(t team.Team) (gormodel.Team, error) {
	var team gormodel.Team

	err := r.db.
		Where("public_id = ?", t.PublicID()).
		First(&team).
		Error
	if err != nil {
		return team, nil
	}

	// TODO: implement members update
	team.Name = t.Name()

	err = r.db.Save(&team).Error
	if err != nil {
		return gormodel.Team{}, err
	}

	return team, nil
}

func (r teamRepo) Delete(publicID string) error {
	var team gormodel.Team

	err := r.db.
		Where("public_id = ?", publicID).
		First(&team).
		Delete(&team).
		Error
	if err != nil {
		return nil
	}

	return nil
}

func (r teamRepo) GetByID(publicID string) (gormodel.Team, error) {
	var team gormodel.Team

	err := r.db.
		Where("public_id = ?", publicID).
		First(&team).
		Error
	if err != nil {
		return team, nil
	}

	return team, nil
}