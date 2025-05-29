package main

import (
	"context"
	"encoding/xml"
	"html"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (r *RSSFeed) unescape() {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)

	for i, v := range r.Channel.Item {
		r.Channel.Item[i].Title = html.UnescapeString(v.Title)
		r.Channel.Item[i].Description = html.UnescapeString(v.Description)
	}
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

	feed.unescape()

	return &feed, nil
}
