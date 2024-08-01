package team_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/team"
	"github.com/bernardinorafael/kn-server/internal/core/domain/entity/user"
	"github.com/stretchr/testify/assert"
)

func TestTeam_New(t *testing.T) {
	t.Run("Should create a new team", func(t *testing.T) {
		u, _ := user.New(
			"john doe",
			"john.doe@email.com",
			".John123",
			"75838249072",
			"48988781289",
			nil,
		)

		tm, err := team.New(u.PublicID, "John's team")
		assert.Nil(t, err)

		assert.Equal(t, tm.OwnerID(), u.PublicID)
		assert.Equal(t, tm.Name(), "John's team")
	})
}

func TestTeam_AddMembers(t *testing.T) {
	t.Run("Should add new members to a team", func(t *testing.T) {
		// jane is the team owner
		jane, _ := user.New(
			"jane doe",
			"jane.doe@email.com",
			".Jane123",
			"75838249072",
			"48988781289",
			nil,
		)

		// bob is the team member
		bob, _ := user.New(
			"bob steel",
			"bob.steel@email.com",
			".Bob123",
			"75838249072",
			"48988781289",
			nil,
		)

		tm, err := team.New(jane.PublicID, "Jane's team")
		assert.Nil(t, err)

		err = tm.AddMembers(*bob)

		assert.NotNil(t, err)
		assert.Len(t, tm.Members(), 1)
	})

	t.Run("Should throw an error if the owner is added as member", func(t *testing.T) {
		jane, _ := user.New(
			"jane doe",
			"jane.doe@email.com",
			".Jane123",
			"75838249072",
			"48988781289",
			nil,
		)

		tm, _ := team.New(jane.PublicID, "Jane's team")

		members := []user.User{*jane}
		err := tm.AddMembers(members...)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "the owner cannot be added as a member")
	})
}
