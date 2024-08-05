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

	var sb strings.Builder
	for _, char := range s.value {
		if char == ' ' {
			sb.WriteRune(char)
			continue
		}
		switch char {
		case 'á', 'à', 'ã', 'â', 'ä':
			sb.WriteRune('a')
		case 'é', 'è', 'ê', 'ë':
			sb.WriteRune('e')
		case 'í', 'ì', 'î', 'ï':
			sb.WriteRune('i')
		case 'ó', 'ò', 'õ', 'ô', 'ö':
			sb.WriteRune('o')
		case 'ú', 'ù', 'û', 'ü':
			sb.WriteRune('u')
		case 'ç':
			sb.WriteRune('c')
		default:
			sb.WriteRune(char)
		}
	}
	s.value = sb.String()

	r := regexp.MustCompile(`[^a-z0-9]+`)
	s.value = r.ReplaceAllString(s.value, "-")

	r, _ = regexp.Compile("-+")
	s.value = r.ReplaceAllString(s.value, "-")

	s.value = strings.Trim(s.value, "-")
}

func (s *slug) Slug() Slug {
	s.slugify()
	return Slug(s.value)
}
