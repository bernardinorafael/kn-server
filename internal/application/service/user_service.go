package service

import (
	"errors"
	"log/slog"
	"time"

	"github.com/bernardinorafael/gozinho/internal/application/dto"
	"github.com/bernardinorafael/gozinho/internal/application/interfaces"
	"github.com/bernardinorafael/gozinho/internal/domain/entity"
	"github.com/bernardinorafael/gozinho/internal/infra/http/response"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository interfaces.UserRepository
}

func NewUserService(repository interfaces.UserRepository) *UserService {
	return &UserService{repository}
}

func (s *UserService) Save(u *dto.CreateUser) error {
	_user, err := s.repository.GetByEmail(u.Email)
	if err != nil {
		slog.Error("error to find user by email", err)
		return err
	}
	if _user != nil {
		slog.Error("user already taken")
		return errors.New("user already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		slog.Error("failed to encrypt password", err)
		return errors.New("failed to encrypt password")
	}

	user := entity.User{
		ID:         uuid.New().String(),
		Name:       u.Name,
		Username:   u.Username,
		Email:      u.Email,
		PersonalID: u.PersonalID,
		Password:   string(hash),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repository.Save(&user); err != nil {
		slog.Error("error to create user", err)
		return err
	}

	return nil
}

func (s *UserService) GetByID(id string) (*response.UserResponse, error) {
	_user, err := s.repository.GetByID(id)
	if err != nil {
		slog.Error("error to find user by ID", err)
		return nil, err
	}
	if _user == nil {
		slog.Error("user not found")
		return nil, errors.New("user not found")
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

func (s *UserService) Update(u *dto.UpdateUser, id string) error {
	user, err := s.repository.GetByID(id)
	if err != nil {
		slog.Error("error to find user by ID", err)
		return err
	}
	if user == nil {
		slog.Error("user not found")
		return errors.New("user not found")
	}

	if u.Email != "" {
		user, err := s.repository.GetByEmail(u.Email)
		if err != nil {
			slog.Error("error to find user by e-mail")
			return err
		}
		if user != nil {
			slog.Error("user already exists")
			return errors.New("e-mail already taken")
		}
	}

	updated := entity.User{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
	}

	if err := s.repository.Update(&updated); err != nil {
		slog.Error("error to update user", err)
		return err
	}

	return nil
}

func (s *UserService) Delete(id string) error {
	user, err := s.repository.GetByID(id)
	if err != nil {
		slog.Error("error finding user by ID", user)
		return err
	}
	if user == nil {
		slog.Error("user not found")
		return errors.New("user not found")
	}

	if err := s.repository.Delete(id); err != nil {
		slog.Error("error to delete user", err)
		return err
	}

	return nil
}

func (s *UserService) GetAll() (*response.AllUsersResponse, error) {
	_users, err := s.repository.GetAll()
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

func (s *UserService) UpdatePassword(u *dto.UpdatePassword, id string) error {
	_user, err := s.repository.GetByID(id)
	if err != nil {
		slog.Error("error do find user by ID", err)
		return err
	}
	if _user == nil {
		slog.Error("user not found")
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(_user.Password), []byte(u.PreviousPassword))
	if err != nil {
		slog.Error("invalid previous password")
		return errors.New("invalid previous password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(_user.Password), []byte(u.PreviousPassword))
	if err == nil {
		slog.Error("new password is equal to previous password")
		return errors.New("new password is equal to previous password")
	}

	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		slog.Error("error to encrypt password", err)
		return err
	}

	if err := s.repository.UpdatePassword(string(encrypted), id); err != nil {
		slog.Error("failed to update password", err)
		return err
	}

	return nil
}
