package slug

import (
	"errors"
	"regexp"
	"strings"
)

type Slug string

type slug struct {
	value string
}

func New(name string) (*slug, error) {
	if len(name) == 0 {
		return nil, errors.New("invalid entrypoint slug")
	}

	slug := slug{value: name}
	return &slug, nil
}

func (s *slug) slugify() {
	s.value = strings.ToLower(s.value)

	var res strings.Builder
	for _, char := range s.value {
		if char == ' ' {
			res.WriteRune(char)
			continue
		}
		switch char {
		case 'á', 'à', 'ã', 'â', 'ä':
			res.WriteRune('a')
		case 'é', 'è', 'ê', 'ë':
			res.WriteRune('e')
		case 'í', 'ì', 'î', 'ï':
			res.WriteRune('i')
		case 'ó', 'ò', 'õ', 'ô', 'ö':
			res.WriteRune('o')
		case 'ú', 'ù', 'û', 'ü':
			res.WriteRune('u')
		case 'ç':
			res.WriteRune('c')
		default:
			res.WriteRune(char)
		}
	}
	s.value = res.String()

	r := regexp.MustCompile(`[^a-z0-9]+`)
	s.value = r.ReplaceAllString(s.value, "-")

	r, _ = regexp.Compile("-+")
	s.value = r.ReplaceAllString(s.value, "-")

	s.value = strings.Trim(s.value, "-")
}

func (s *slug) GetSlug() Slug {
	s.slugify()
	return Slug(s.value)
}
