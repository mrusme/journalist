package rss

import (
	// log "github.com/sirupsen/logrus"
	"context"
	"time"

	"github.com/mmcdole/gofeed"
)

type Client struct {
  parser        *gofeed.Parser
  url           string
  Feed          *gofeed.Feed
  Items         *[]*gofeed.Item
  UpdatedAt     time.Time
}

func NewClient(feedUrl string) (*Client, error) {
  client := new(Client)
  client.parser = gofeed.NewParser()
  client.url = feedUrl

  if err := client.Sync(); err != nil {
    return nil, err
  }

  return client, nil
}

func (c *Client) Sync() (error) {
  ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
  defer cancel()

  feed, err := c.parser.ParseURLWithContext(c.url, ctx)
  if err != nil {
    return err
  }

  c.Feed = feed
  c.Items = &feed.Items
  c.UpdatedAt = time.Now()

  return nil
}

