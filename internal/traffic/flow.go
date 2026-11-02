package traffic

import (
	"errors"
	"net/http"
	"time"
)

const (
	DefaultCapacity = 500
	MaxBodySize     = 1 << 20
	FlowEventName   = "traffic:flow"
)

var ErrNotFound = errors.New("flow not found")

type Summary struct {
	ID         string    `json:"id"`
	Method     string    `json:"method"`
	URL        string    `json:"url"`
	Status     int       `json:"status"`
	StartTime  time.Time `json:"startTime"`
	DurationMs int64     `json:"durationMs"`
	Size       int       `json:"size"`
	Error      string    `json:"error,omitempty"`
}

type Message struct {
	Header    http.Header `json:"header"`
	Body      string      `json:"body"`
	Truncated bool        `json:"truncated"`
	Streamed  bool        `json:"streamed"`
}

type Flow struct {
	Summary  Summary  `json:"summary"`
	Request  Message  `json:"request"`
	Response *Message `json:"response,omitempty"`
}

// shortcut: body is stored as a string, so binary bodies show garbled; add base64 + content-type sniffing if the UI needs them.
func NewMessage(
	header http.Header,
	body []byte,
	streamed bool,
) Message {
	m := Message{
		Header:   header.Clone(),
		Streamed: streamed,
	}
	if len(body) > MaxBodySize {
		body = body[:MaxBodySize]
		m.Truncated = true
	}

	m.Body = string(body)

	return m
}
