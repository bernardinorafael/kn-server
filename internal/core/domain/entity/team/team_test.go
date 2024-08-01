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

		tm, err := team.New(u.PublicID, "My team")
		assert.Nil(t, err)

		assert.Equal(t, tm.OwnerID(), u.PublicID)
		assert.Equal(t, tm.Name(), "My team")
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

		tm, err := team.New(jane.PublicID, "My team")
		assert.Nil(t, err)

		members := []user.User{*bob}
		tm.AddMembers(members...)

		assert.Len(t, tm.Members(), 1)
	})
}
