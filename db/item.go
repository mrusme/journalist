package db

import (
  "time"
)

type Item struct {
  ID                string          `json:"-"`
  IncID             int             `json:"inc_id,omitempty"`
  Title             string          `json:"title,omitempty"`
  URL               string          `json:"url,omitempty"`
  Author            string          `json:"author,omitempty"`
  RssBody           string          `json:"rss_body,omitempty"`
  SiteBody          string          `json:"site_body,omitempty"`
  SiteBodyOptimized string          `json:"site_body_optimized,omitempty"`
  IsRead            bool            `json:"is_read,omitempty"`
  IsSaved           bool            `json:"is_saved,omitempty"`
  FaviconData       string          `json:"favicon_data,omitempty"`
  CreatedAt         time.Time       `json:"created_at,omitempty"`
}

func (item *Item) SetIDFromDatabaseKey(key string) (error) {
  var err error
  item.ID, err = GetIDFromDatabaseKey(key)
  return err
}
