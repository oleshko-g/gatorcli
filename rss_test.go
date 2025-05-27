package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchFeed(t *testing.T) {
	const wagsLaneRSS = "https://www.wagslane.dev/index.xml"
	ctx := context.Background()

	rss, err := fetchFeed(ctx, wagsLaneRSS)
	require.NoError(t, err)
	require.NotEmpty(t, rss)
}
