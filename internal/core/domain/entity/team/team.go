package team

import (
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"

	"github.com/google/uuid"
)

type Team struct {
	publicID  string
	owner     user.User
	name      string
	members   []user.User
	createdAt time.Time
}

func New(owner user.User, name string, members ...user.User) (*Team, error) {
	team := &Team{
		publicID:  uuid.NewString(),
		owner:     owner,
		name:      name,
		members:   members,
		createdAt: time.Now(),
	}

	if err := team.validate(); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *Team) validate() error {
	if t.owner.PublicID == "" {
		return errors.New("the team must have an owner")
	}

	if t.name == "" {
		return errors.New("the team name cannot be empty")
	}

	status := t.owner.Enabled
	if status == false {
		return errors.New("only activated users can create a team")
	}

	return nil
}

func (t *Team) AddMembers(members ...user.User) error {
	for _, m := range members {
		if m.PublicID == t.owner.PublicID {
			return errors.New("the owner cannot be added as a member")
		}
	}
	t.members = append(t.members, members...)
	return nil
}

func (t *Team) PublicID() string     { return t.publicID }
func (t *Team) Owner() user.User     { return t.owner }
func (t *Team) Name() string         { return t.name }
func (t *Team) Members() []user.User { return t.members }
func (t *Team) CreatedAt() time.Time { return t.createdAt }
