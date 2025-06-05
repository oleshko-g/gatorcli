package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/oleshko-g/gatorcli/internal/database"
	"github.com/stretchr/testify/require"
)

func TestFetchFeed(t *testing.T) {
	const wagsLaneRSS = "https://www.wagslane.dev/index.xml"
	ctx := context.Background()

	rss, err := fetchFeed(ctx, wagsLaneRSS)
	require.NoError(t, err)
	require.NotEmpty(t, rss)
}

func TestScrapeFeeds(t *testing.T) {
	s := setState()
	db, errDB := openPostgresDB(s.cfg.DataBaseURL)
	if errDB != nil {
		fmt.Fprintln(os.Stderr, errDB)
		os.Exit(1)
	}
	s.db = database.New(db)

	err := scrapeFeeds(&s)
	require.NoError(t, err)
}
