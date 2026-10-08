package models

import "time"

type TagChecks struct {
	State       string    `json:"state"`
	StartedAt   time.Time `json:"startedAt,omitzero"`
	CompletedAt time.Time `json:"completedAt,omitzero"`
}

func (c TagChecks) IsRunning() bool {
	return (c.State == "PENDING" || c.State == "EXPECTED") && !c.StartedAt.IsZero()
}

func (c TagChecks) Elapsed(now time.Time) (time.Duration, bool) {
	switch {
	case c.StartedAt.IsZero():
		return 0, false
	case !c.CompletedAt.IsZero():
		return c.CompletedAt.Sub(c.StartedAt), true
	case c.IsRunning():
		return now.Sub(c.StartedAt), true
	default:
		return 0, false
	}
}
