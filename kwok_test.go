package kwok_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	kwok "github.com/pablochacin/testcontainers-go-kwok"
)

func TestKwok(t *testing.T) {
	ctx := context.Background()

	ctr, err := kwok.Run(ctx, "ghcr.io/pablochacin/kwok:latest", kwok.WithNodes(2))
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)
}
