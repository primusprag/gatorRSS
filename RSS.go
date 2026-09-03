package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Item        []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error with request: %w", err)
	}
	defer res.Body.Close()

	req.Header.Set("User-Agent", "gator")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body %w", err)
	}

	var RSSPtr RSSFeed
	err = xml.Unmarshal(body, &RSSPtr)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling feed %w", err)
	}

	/* Given solutoin:
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}
	*/

	return htmlDecode(&RSSPtr), nil
}

func htmlDecode(oldRSS *RSSFeed) *RSSFeed {
	var newItems []RSSItem

	for _, item := range oldRSS.Channel.Item {
		newItems = append(newItems, RSSItem{
			Title:       html.UnescapeString(item.Title),
			Link:        item.Link,
			Description: html.UnescapeString(item.Description),
			PubDate:     item.PubDate,
		})
	}

	newRSS := &RSSFeed{
		Channel: RSSChannel{
			Title:       html.UnescapeString(oldRSS.Channel.Title),
			Link:        oldRSS.Channel.Link,
			Description: html.UnescapeString(oldRSS.Channel.Description),
			Item:        newItems,
		},
	}

	return newRSS
}
