package service

import (
	"database/sql"
	"errors"
	"log/slog"
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

func (s *accountService) Save(u *dto.UserInput) error {
	_, err := s.svc.ar.GetByEmail(u.Email)
	if err == nil {
		slog.Error("email already taken")
		return errors.New("email already taken")
	}

	password, err := crypto.HashPassword(u.Password)
	if err != nil {
		slog.Error("failed to encrypt password", err)
		return errors.New("failed to encrypt password")
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

	if err := s.svc.ar.Save(&user); err != nil {
		slog.Error("error to create user", err)
		return err
	}

	return nil
}

func (s *accountService) GetByID(id string) (*response.UserResponse, error) {
	_user, err := s.svc.ar.GetByID(id)
	if err != nil {
		slog.Error("user not found", err, slog.String("pkg", "service"))
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

func (s *accountService) Update(u *dto.UpdateUser, id string) error {
	user, err := s.svc.ar.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("error to find user by ID", err)
			return err
		}
	}
	if user == nil {
		slog.Error("user not found")
		return errors.New("user not found")
	}

	if u.Email != "" {
		user, err := s.svc.ar.GetByEmail(u.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				slog.Error("error to find user by e-mail")
				return err
			}
		}
		if user == nil {
			slog.Error("user already exists")
			return errors.New("e-mail already taken")
		}
	}

	updated := entity.Account{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
	}

	if err := s.svc.ar.Update(&updated); err != nil {
		slog.Error("error to update user", err)
		return err
	}

	return nil
}

func (s *accountService) Delete(id string) error {
	user, err := s.svc.ar.GetByID(id)
	if err != nil {
		slog.Error("error finding user by ID", user)
		return err
	}
	if user == nil {
		slog.Error("user not found")
		return errors.New("user not found")
	}

	if err := s.svc.ar.Delete(id); err != nil {
		slog.Error("error to delete user", err)
		return err
	}

	return nil
}

func (s *accountService) GetAll() (*response.AllUsersResponse, error) {
	_users, err := s.svc.ar.GetAll()
	if err != nil {
		slog.Error("error to find many users", err)
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

func (s *accountService) UpdatePassword(u *dto.UpdatePassword, id string) error {
	user, err := s.svc.ar.GetByID(id)
	if err != nil {
		slog.Error("error do find user by ID", err, slog.String("pkg", "service"))
		return err
	}

	err = crypto.CheckPassword(user.Password, u.PreviousPassword)
	if err != nil {
		slog.Error("invalid previous password", slog.String("pkg", "service"))
		return errors.New("invalid previous password")
	}

	password, err := crypto.HashPassword(u.Password)
	if err != nil {
		slog.Error("error to encrypt password", err, slog.String("pkg", "service"))
		return err
	}

	if err := s.svc.ar.UpdatePassword(password, id); err != nil {
		slog.Error("failed to update password", err, slog.String("pkg", "service"))
		return err
	}

	return nil
}
