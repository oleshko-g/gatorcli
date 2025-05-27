package main

import (
	"context"
	"encoding/xml"
	"net/http"
)

type RSSFeed = struct {
	Channel struct {
		Title       string  `xml:"title"`
		Link        string  `xml:"link"`
		Description string  `xml:"description"`
		Item        RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, errReq := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if errReq != nil {
		return nil, errReq
	}
	req.Header.Set("user-agent", "gatorcli")

	c := http.Client{}

	res, errRes := c.Do(req)
	if errRes != nil {
		return nil, errRes
	}
	defer res.Body.Close()

	var feed RSSFeed
	d := xml.NewDecoder(res.Body)
	errDecode := d.Decode(&feed)
	if errDecode != nil {
		return nil, errDecode
	}

	return &feed, nil
}
