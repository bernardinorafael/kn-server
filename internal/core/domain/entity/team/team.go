package team

import (
	"errors"
	"time"

	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"

	"github.com/google/uuid"
)

type Team struct {
	publicID  string
	ownerID   string
	name      string
	members   []user.User
	createdAt time.Time
}

func New(ownerID, name string) (*Team, error) {
	team := &Team{
		publicID:  uuid.NewString(),
		ownerID:   ownerID,
		name:      name,
		members:   nil,
		createdAt: time.Now(),
	}

	if err := team.validate(); err != nil {
		return nil, err
	}

	return team, nil
}

func (t *Team) validate() error {
	if t.ownerID == "" {
		return errors.New("the team must have an owner")
	}

	if t.name == "" {
		return errors.New("the team name cannot be empty")
	}

	return nil
}

func (t *Team) AddMember(member user.User) error {
	if member.PublicID() == t.ownerID {
		return errors.New("the owner cannot be added as a member")
	}

	for _, m := range t.Members() {
		if m.PublicID() == member.PublicID() {
			return errors.New("a member cannot be added twice")
		}
	}
	t.members = append(t.members, member)

	return nil
}

func (t *Team) PublicID() string     { return t.publicID }
func (t *Team) OwnerID() string      { return t.ownerID }
func (t *Team) Name() string         { return t.name }
func (t *Team) Members() []user.User { return t.members }
func (t *Team) CreatedAt() time.Time { return t.createdAt }
