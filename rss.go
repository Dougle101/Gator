package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/Dougle101/Gator/internal/database"
	"github.com/google/uuid"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request: %v\n", err)
	}
	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending HTTP request: %v\n", err)
	}
	defer res.Body.Close()

	raw_xml, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v\n", err)
	}

	var dec_xml RSSFeed
	err = xml.Unmarshal(raw_xml, &dec_xml)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling XML: %v\n", err)
	}

	dec_xml.Channel.Title = html.UnescapeString(dec_xml.Channel.Title)
	dec_xml.Channel.Description = html.UnescapeString(dec_xml.Channel.Description)
	for i := range dec_xml.Channel.Item {
		item := &dec_xml.Channel.Item[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &dec_xml, nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed.Channel.Title)
	for _, item := range feed.Channel.Item {
		fmt.Printf("- %s\n %s\n", item.Title, item.Link)
	}
	return nil
}

func handlerCreateFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: gator addfeed <Feed Name> <Feed URL>, received %d arguments (expecting 2)", len(cmd.args))
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("Error fetching user: %v\n", err)
	}

	feed_params := database.CreateFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID}

	feed, err := s.db.CreateFeed(context.Background(), feed_params)
	if err != nil {
		return fmt.Errorf("Error creating feed: %v\n", err)
	}
	fmt.Printf("Created Feed: %+v\n", feed)
	return nil
}
