package gormrepo

import (
	"github.com/bernardinorafael/kn-server/internal/core/application/contract"
	"github.com/bernardinorafael/kn-server/internal/infra/database/gorm/gormodel"
	"gorm.io/gorm"
)

/*
 TODO: remove entity/model mapping logic from repositories and do it into service layer
*/

type teamRepo struct{ db *gorm.DB }

func NewTeamRepo(db *gorm.DB) contract.GormTeamRepository {
	return &teamRepo{db}
}

func (r teamRepo) Create(t gormodel.Team) error {
	err := r.db.Create(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func (r teamRepo) Update(t gormodel.Team) (gormodel.Team, error) {
	var team gormodel.Team

	err := r.db.
		Where("public_id = ?", t.PublicID).
		First(&team).
		Error
	if err != nil {
		return team, nil
	}

	// TODO: implement members update
	team.Name = t.Name

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

func (r teamRepo) GetByPublicID(publicID string) (gormodel.Team, error) {
	var team gormodel.Team

	err := r.db.
		Where("public_id = ?", publicID).
		First(&team).
		Error
	if err != nil {
		return team, err
	}

	return team, nil
}
