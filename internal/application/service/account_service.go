package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/bernardinorafael/kn-server/internal/application/contract"
	"github.com/bernardinorafael/kn-server/internal/application/dto"
	"github.com/bernardinorafael/kn-server/internal/domain/entity"
	"github.com/bernardinorafael/kn-server/internal/infra/auth"
	"github.com/bernardinorafael/kn-server/internal/infra/rest/response"
	"github.com/bernardinorafael/kn-server/util/crypto"
	"github.com/google/uuid"
)

type accountService struct {
	service *service
}

func newAccountService(service *service) contract.AccountService {
	return &accountService{service}
}

func (a *accountService) Save(ctx context.Context, u dto.UserInput) error {
	a.service.l.Info(ctx, "process started")
	defer a.service.l.Info(ctx, "process finished")

	exist, _ := a.service.ar.CheckUserExist(u.Email, u.Username, u.Document)
	if exist {
		a.service.l.Error(ctx, "user already taken")
		return errUserAlreadyTaken
	}

	password, err := crypto.EncryptPassword(u.Password)
	if err != nil {
		a.service.l.Errorf(ctx, "error hashing password: %v", err.Error())
		return errHashPassword
	}

	user := entity.Account{
		ID:        uuid.New().String(),
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Password:  password,
		Document:  u.Document,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.service.ar.Save(&user); err != nil {
		a.service.l.Errorf(ctx, "error creating user: %s", err.Error())
		return errCreateUser
	}
	return nil
}

func (a *accountService) GetByID(ctx context.Context, id string) (*response.UserResponse, error) {
	a.service.l.Info(ctx, "process started")
	defer a.service.l.Info(ctx, "process finished")

	u, err := a.service.ar.GetByID(id)
	if err != nil {
		a.service.l.Errorf(ctx, "error to find user: %s", err.Error())
		return nil, errUserNotFound
	}
	userInDB := u

	user := response.UserResponse{
		ID:        userInDB.ID,
		Email:     userInDB.Email,
		Document:  userInDB.Document,
		Name:      userInDB.Name,
		Username:  userInDB.Username,
		CreatedAt: userInDB.CreatedAt,
	}
	return &user, nil
}

func (a *accountService) Update(ctx context.Context, u dto.UpdateUser, id string) error {
	a.service.l.Info(ctx, "process started")
	defer a.service.l.Info(ctx, "process finished")

	_, err := a.service.ar.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.service.l.Errorf(ctx, "error to find user by ID: %s", err.Error())
			return errUserNotFound
		}
	}

	if u.Email != "" {
		_, err := a.service.ar.GetByEmail(u.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				a.service.l.Errorf(ctx, "unable to find user by email: %s", err.Error())
				return errEmailNotFound
			}
		}
	}

	account := entity.Account{
		Name:     u.Name,
		Username: u.Username,
		Email:    u.Email,
	}

	if err := a.service.ar.Update(&account); err != nil {
		a.service.l.Errorf(ctx, "error update user: %s", err.Error())
		return errUpdateUser
	}
	return nil
}

func (a *accountService) Delete(ctx context.Context, id string) error {
	a.service.l.Info(ctx, "process started")
	defer a.service.l.Info(ctx, "process finished")

	_, err := a.service.ar.GetByID(id)
	if err != nil {
		a.service.l.Errorf(ctx, "error find user by ID: %s", err.Error())
		return errUserNotFound
	}

	if err := a.service.ar.Delete(id); err != nil {
		a.service.l.Errorf(ctx, "error deleting user: %s", err.Error())
		return errDeleteUser
	}
	return nil
}

func (a *accountService) GetAll(ctx context.Context) (*response.AllUsersResponse, error) {
	a.service.l.Info(ctx, "process started")
	defer a.service.l.Info(ctx, "process finished")

	u, err := a.service.ar.GetAll()
	if err != nil {
		a.service.l.Errorf(ctx, "error find users: %s", err.Error())
		return nil, errGetManyUsers
	}
	_users := u

	users := response.AllUsersResponse{}
	for _, user := range _users {
		usersResponse := response.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Document:  user.Document,
			CreatedAt: user.CreatedAt,
		}
		users.Users = append(users.Users, usersResponse)
	}
	return &users, nil
}

func (a *accountService) UpdatePassword(ctx context.Context, u dto.UpdatePassword, id string) error {
	a.service.l.Info(ctx, "process started")
	defer a.service.l.Info(ctx, "process finished")

	oldPassDB, err := a.service.ar.GetPassword(id)
	if err != nil {
		a.service.l.Errorf(ctx, "error find user by ID: %s", err.Error())
		return errUserNotFound
	}

	if err = crypto.CheckPassword(oldPassDB, u.OldPassword); err != nil {
		a.service.l.Errorf(ctx, "invalid old password: %s", err.Error())
		return errInvAlidCredential
	}

	err = crypto.CheckPassword(oldPassDB, u.Password)
	if err == nil {
		a.service.l.Error(ctx, "old password is equal to current one")
		return errEqualPasswords
	}

	p, err := crypto.EncryptPassword(u.Password)
	if err != nil {
		a.service.l.Errorf(ctx, "error hashing password: %s", err.Error())
		return errHashPassword
	}
	password := p

	err = a.service.ar.UpdatePassword(password, id)
	if err != nil {
		a.service.l.Errorf(ctx, "error updating password: %s", err.Error())
		return errUpdateUser

	}
	return nil
}

func (a *accountService) Login(ctx context.Context, input dto.Login) (*entity.Account, error) {
	a.service.l.Info(ctx, "process started")
	defer a.service.l.Info(ctx, "process finished")

	acc, err := a.service.ar.GetByEmail(input.Email)
	if err != nil {
		a.service.l.Errorf(ctx, "cannot find user by email: %s", err.Error())
		return nil, errEmailNotFound
	}

	ctx = context.WithValue(ctx, auth.UserIDKey, acc.ID)

	_, err = a.service.ar.GetPassword(acc.ID)
	if err != nil {
		a.service.l.Errorf(ctx, "error to get user password: %s", err.Error())
		return nil, errInvAlidCredential
	}

	err = crypto.CheckPassword(acc.Password, input.Password)
	if err != nil {
		a.service.l.Errorf(ctx, "password does not match: %s", err.Error())
		return nil, errInvAlidCredential
	}

	a.service.l.Infow(ctx, "user info used to login",
		slog.Group("user-info"),
		slog.String("user-id", acc.ID),
		slog.String("user-name", acc.Name),
	)

	account := &entity.Account{
		ID:        acc.ID,
		Name:      acc.Name,
		Username:  acc.Username,
		Email:     acc.Email,
		Password:  acc.Password,
		Document:  acc.Document,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return account, nil
}
