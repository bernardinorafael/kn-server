package team

import (
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"

	"github.com/google/uuid"
)

type Team struct {
	publicId  string
	ownerId   string
	name      string
	members   []user.User
	createdAt time.Time
}

func New(ownerId, name string, members ...user.User) (*Team, error) {
	team := &Team{
		publicId:  uuid.NewString(),
		ownerId:   ownerId,
		name:      name,
		members:   members,
		createdAt: time.Now(),
	}

	if err := team.validate(); err != nil {
		return nil, err
	}

	return team, nil
}

func (o *Team) validate() error {
	if o.ownerId == "" {
		return errors.New("ownerId cannot be empty")
	}
	if o.name == "" {
		return errors.New("name cannot be empty")
	}

	return nil
}

func (o *Team) AddMembers(members ...user.User) error {
	for _, m := range members {
		if m.PublicID == o.ownerId {
			return errors.New("the owner cannot be added as a member")
		}
	}
	o.members = append(o.members, members...)
	return nil
}

func (o *Team) PublicID() string     { return o.publicId }
func (o *Team) OwnerID() string      { return o.ownerId }
func (o *Team) Name() string         { return o.name }
func (o *Team) Members() []user.User { return o.members }
func (o *Team) CreatedAt() time.Time { return o.createdAt }
