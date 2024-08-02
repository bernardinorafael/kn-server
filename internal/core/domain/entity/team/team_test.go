package team_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/team"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTeam_New(t *testing.T) {
	t.Run("Should create a new team", func(t *testing.T) {
		john, _ := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "john doe",
			Email:    "john.doe@email.com",
			Password: ".John1234",
			Document: "75838249072",
			Phone:    "48988781289",
			TeamID:   nil,
		})

		tm, err := team.New(john.PublicID(), "John's team")

		assert.Nil(t, err)
		assert.Equal(t, tm.OwnerID(), john.PublicID())
		assert.Equal(t, tm.Name(), "John's team")
	})
}

func TestTeam_AddMembers(t *testing.T) {
	t.Run("Should add new members to a team", func(t *testing.T) {
		// jane is the owner
		jane, _ := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "jane doe",
			Email:    "jane.doe@email.com",
			Password: ".Jane1234",
			Document: "75838249072",
			Phone:    "48988781289",
			TeamID:   nil,
		})

		// bob is the team member
		bob, _ := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "bob doe",
			Email:    "bob.doe@email.com",
			Password: ".Bob1234",
			Document: "75838249072",
			Phone:    "48988781289",
			TeamID:   nil,
		})

		tm, err := team.New(jane.PublicID(), "Jane's team")
		err = tm.AddMember(*bob)

		assert.NotNil(t, err)
		assert.Len(t, tm.Members(), 1)
	})

	t.Run("Should throw an error if the owner is added as member", func(t *testing.T) {
		jane, _ := user.New(user.Params{
			PublicID: uuid.NewString(),
			Name:     "jane doe",
			Email:    "jane.doe@email.com",
			Password: ".Jane1234",
			Document: "75838249072",
			Phone:    "48988781289",
			TeamID:   nil,
		})

		tm, _ := team.New(jane.PublicID(), "Jane's team")
		err := tm.AddMember(*jane)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "the owner cannot be added as a member")
	})
}
