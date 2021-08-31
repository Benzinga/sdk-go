package bztime

import (
	"errors"
	"strings"
	"time"
)

const RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"

type Time struct {
	time.Time
}

func FromTime(t time.Time) Time {
	n := Time{}
	n.Time = t

	return n
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (t Time) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("bztime.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(RFC3339Milli)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, RFC3339Milli)
	b = append(b, '"')

	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *Time) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse(RFC3339Milli, strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}

	t.Time = tt

	return nil
}
