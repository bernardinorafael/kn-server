package slug_test

import (
	"testing"

	"github.com/bernardinorafael/kn-server/internal/domain/valueobj/slug"
	"github.com/stretchr/testify/assert"
)

func TestSlugifyString(t *testing.T) {
	t.Run("Should slugify a string with special characters", func(t *testing.T) {
		s := "EsTe é Um téste da Função SLUGIFY"
		slug, err := slug.New(s)

		assert.Nil(t, err)
		assert.Equal(t, string(slug.GetSlug()), "este-e-um-teste-da-funcao-slugify")
	})
}
