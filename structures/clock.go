package structures

import "time"

type Clock struct {
	Ex int64
}

func (c *Clock) Now() time.Time          { return time.Now() }
func (c *Clock) Duration() time.Duration { return time.Duration(c.Ex) * time.Minute }
