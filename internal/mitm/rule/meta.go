package rule

import (
	"errors"
	"time"
)

const (
	MaxRules       = 5
	DefaultTimeout = 500 * time.Millisecond
)

var (
	ErrMaxRules      = errors.New("rule: max rules reached")
	ErrNotFound      = errors.New("rule: not found")
	ErrExists        = errors.New("rule: already exists")
	ErrScriptTimeout = errors.New("rule: script timeout")
	ErrInvalidRule   = errors.New("rule: invalid rule")
)

type Rule struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Script  string `json:"script"`
}
