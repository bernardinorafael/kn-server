package team

import (
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"

	"github.com/google/uuid"
)

var (
	ErrEmptyOwnerID = errors.New("ownerId cannot be empty")
	ErrEmptyOrgName = errors.New("name cannot be empty")
)

type Team struct {
	publicId  string
	ownerId   string
	name      string
	members   []user.User
	createdAt time.Time
}

func New(ownerId, name string, members ...user.User) (*Team, error) {
	org := &Team{
		publicId:  uuid.NewString(),
		ownerId:   ownerId,
		name:      name,
		members:   members,
		createdAt: time.Now(),
	}

	if err := org.validate(); err != nil {
		return nil, err
	}

	return org, nil
}

func (o *Team) validate() error {
	if o.ownerId == "" {
		return ErrEmptyOwnerID
	}
	if o.name == "" {
		return ErrEmptyOrgName
	}
	return nil
}

func (o *Team) AddMembers(members ...user.User) {
	o.members = append(o.members, members...)
}

func (o *Team) PublicID() string     { return o.publicId }
func (o *Team) OwnerID() string      { return o.ownerId }
func (o *Team) Name() string         { return o.name }
func (o *Team) Members() []user.User { return o.members }
func (o *Team) CreatedAt() time.Time { return o.createdAt }
