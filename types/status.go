package types

import (
	"fmt"
	"strings"
	"time"
)

type StatusEntry struct {
	Timestamp time.Time `json:"time"`
	Message   string    `json:"message"`
}

type Status struct {
	Success  bool          `json:"status"`
	Messages []StatusEntry `json:"messages"`
}

func NewStatus() Status {
	return Status{Success: true}
}

func (s *Status) LogF(format string, a ...interface{}) {
	s.Log(fmt.Sprintf(format, a...))
}

func (s *Status) Fail(err error) {
	s.Success = false
	s.Log(err.Error())
}

func (s *Status) Log(message ...string) {
	s.Messages = append(s.Messages, StatusEntry{
		Timestamp: time.Now(),
		Message:   strings.Join(message, " "),
	})
}
