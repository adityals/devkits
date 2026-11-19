package convert

import (
	"errors"
	"fmt"
)

type Lang string

const (
	LangGo Lang = "go"
	LangTS Lang = "ts"
)

type Route struct {
	From Lang
	To   Lang
}

// Converter turns source text of one language into another.
type Converter interface {
	Convert(src string) (string, error)
}

var ErrUnsupportedRoute = errors.New("convert: unsupported route")

type convertService struct {
	routes map[Route]Converter
}

func NewConvertService(routes map[Route]Converter) *convertService {
	return &convertService{routes: routes}
}

func (s *convertService) Convert(
	from Lang,
	to Lang,
	input string,
) (string, error) {
	c, ok := s.routes[Route{From: from, To: to}]
	if !ok {
		return "", fmt.Errorf("%s → %s: %w", from, to, ErrUnsupportedRoute)
	}

	out, err := c.Convert(input)
	if err != nil {
		return "", fmt.Errorf("%s → %s: %w", from, to, err)
	}

	return out, nil
}
