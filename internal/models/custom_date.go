package models

import (
	"fmt"
	"strings"
	"time"
)

type CustomDate struct {
	time.Time
}

const format = "02.01.2006"

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	if s == "null" {
		return
	}
	c.Time, err = time.Parse(format, s)
	return
}

func (c CustomDate) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(format))), nil
}
