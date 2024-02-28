package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/bernardinorafael/gozinho/internal/application/contract"
	"github.com/bernardinorafael/gozinho/internal/application/dto"
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
	"github.com/bernardinorafael/gozinho/internal/infra/rest/response"
	"github.com/bernardinorafael/gozinho/util/crypto"
	"github.com/google/uuid"
)

type accountService struct {
	svc *service
}

func newAccountService(svc *service) contract.AccountService {
	return &accountService{
		svc: svc,
	}
}

func (a *accountService) Save(ctx context.Context, u *dto.UserInput) error {
	a.svc.l.Info(ctx, "process started")
	defer a.svc.l.Info(ctx, "process finished")

	exist, err := a.svc.ar.CheckUserExist(u.Email, u.Username, u.PersonalID)
	if exist {
		a.svc.l.Error(ctx, "user already taken")
		return errors.New("user already taken")
	}

	password, err := crypto.EncryptPassword(u.Password)
	if err != nil {
		a.svc.l.Errorf(ctx, "error hashing password: %v", err.Error())
		return err
	}

	user := entity.Account{
		ID:         uuid.New().String(),
		Name:       u.Name,
		Username:   u.Username,
		Email:      u.Email,
		Password:   password,
		PersonalID: u.PersonalID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := a.svc.ar.Save(&user); err != nil {
		a.svc.l.Errorf(ctx, "error creating user: %s", err.Error())
		return err
	}

	return nil
}

func (a *accountService) GetByID(ctx context.Context, id string) (*response.UserResponse, error) {
	a.svc.l.Info(ctx, "process started")
	defer a.svc.l.Info(ctx, "process finished")

	_user, err := a.svc.ar.GetByID(id)
	if err != nil {
		a.svc.l.Errorf(ctx, "error to find user: %s", err.Error())
		return nil, err
	}

	user := response.UserResponse{
		ID:         _user.ID,
		Email:      _user.Email,
		PersonalID: _user.PersonalID,
		Name:       _user.Name,
		Username:   _user.Username,
		Active:     _user.Active,
		CreatedAt:  _user.CreatedAt,
	}

	return &user, nil
}

func (a *accountService) Update(ctx context.Context, u *dto.UpdateUser, id string) error {
	a.svc.l.Info(ctx, "process started")
	defer a.svc.l.Info(ctx, "process finished")

	_, err := a.svc.ar.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.svc.l.Errorf(ctx, "error to find user by ID: %s", err.Error())
			return err
		}
	}

	if u.Email != "" {
		_, err := a.svc.ar.GetByEmail(u.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				a.svc.l.Errorf(ctx, "unable to find user by email: %s", err.Error())
				return err
			}
		}
	}

	updated := entity.Account{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
	}

	if err := a.svc.ar.Update(&updated); err != nil {
		a.svc.l.Errorf(ctx, "error update user: %s", err.Error())
		return err
	}

	return nil
}

func (a *accountService) Delete(ctx context.Context, id string) error {
	a.svc.l.Info(ctx, "process started")
	defer a.svc.l.Info(ctx, "process finished")

	_, err := a.svc.ar.GetByID(id)
	if err != nil {
		a.svc.l.Errorf(ctx, "error find user by ID: %s", err.Error())
		return err
	}

	if err := a.svc.ar.Delete(id); err != nil {
		a.svc.l.Errorf(ctx, "error deleting user: %s", err.Error())
		return err
	}

	return nil
}

func (a *accountService) GetAll(ctx context.Context) (*response.AllUsersResponse, error) {
	a.svc.l.Info(ctx, "process started")
	defer a.svc.l.Info(ctx, "process finished")

	_users, err := a.svc.ar.GetAll()
	if err != nil {
		a.svc.l.Errorf(ctx, "error find users: %s", err.Error())
		return nil, err
	}

	users := response.AllUsersResponse{}
	for _, user := range _users {
		usersResponse := response.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Username:   user.Username,
			Email:      user.Email,
			PersonalID: user.PersonalID,
			CreatedAt:  user.CreatedAt,
			Active:     user.Active,
		}
		users.Users = append(users.Users, usersResponse)
	}

	return &users, nil
}

func (a *accountService) UpdatePassword(ctx context.Context, u *dto.UpdatePassword, id string) error {
	a.svc.l.Info(ctx, "process started")
	defer a.svc.l.Info(ctx, "process finished")

	oldPassDB, err := a.svc.ar.GetPassword(id)
	if err != nil {
		a.svc.l.Errorf(ctx, "error find user by ID: %s", err.Error())
		return err
	}

	// compare old password sent in req to, old in DB
	err = crypto.CheckPassword(oldPassDB, u.OldPassword)
	if err != nil {
		a.svc.l.Errorf(ctx, "invalid old password: %s", err.Error())
		return err
	}

	// compare old password sent in req, to the new one
	err = crypto.CheckPassword(oldPassDB, u.Password)
	if err == nil {
		a.svc.l.Errorf(ctx, "old password is equal to current one")
		return errors.New("old password is equal to current one")
	}

	password, err := crypto.EncryptPassword(u.Password)
	if err != nil {
		a.svc.l.Errorf(ctx, "error hashing password: %s", err.Error())
		return err
	}

	err = a.svc.ar.UpdatePassword(password, id)
	if err != nil {
		a.svc.l.Errorf(ctx, "error updating password: %s", err.Error())
		return err
	}

	return nil
}
