package kwok_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	kwok "github.com/pablochacin/testcontainers-go-kwok"
	"github.com/testcontainers/testcontainers-go"
)

func TestKwok(t *testing.T) {
	ctx := context.Background()

	ctr, err := kwok.Run(ctx, "ghcr.io/pablochacin/kwok:latest")
	require.NoError(t, err)

	_ , err = ctr.GetKubeConfig(ctx)
	require.NoError(t, err)

	testcontainers.CleanupContainer(t, ctr)
}
