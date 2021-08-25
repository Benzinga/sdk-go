package benzinga

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNews(t *testing.T) {
	client := NewClient(nil).News()
	client.SetAPIToken("68f891dba6ba4c9bae19f4c78c4d64d1")

	ctx := context.Background()

	resp, err := client.Exec(ctx)
	require.NoError(t, err)
	require.Empty(t, resp)
}
