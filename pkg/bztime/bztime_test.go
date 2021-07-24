package bztime

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	j := []byte(`{"time":"2012-02-17T21:01:39.000Z"}`)
	s := "2012-02-17T21:01:39.000Z"

	tm, err := time.Parse(RFC3339Milli, "2012-02-17T21:01:39.000Z")
	require.NoError(t, err)

	t.Run("Test UnmarshalJSON", func(t *testing.T) {
		m := make(map[string]Time)

		require.NoError(t, json.Unmarshal(j, &m))

		require.Equal(t, s, tm.Format(RFC3339Milli))
	})

	t.Run("Test MarshalJSON", func(t *testing.T) {
		tt := Time{}

		tt.Time = tm

		m := map[string]Time{
			"time": tt,
		}

		b, err := json.Marshal(&m)
		require.NoError(t, err)
		require.EqualValues(t, j, b)
	})
}
