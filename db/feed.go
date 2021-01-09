package db

import (
  "time"
)

type Feed struct {
  ID                uint            `json:"id,omitempty"`
  Title             string          `json:"title,omitempty"`
  Description       string          `json:"description,omitempty"`
  Link              string          `json:"link,omitempty"`
  FeedLink          string          `json:"feed_link,omitempty"`
  Author            string          `json:"author,omitempty"`
  Language          string          `json:"language,omitempty"`
  Image             string          `json:"image,omitempty"`
  Copyright         string          `json:"copyright,omitempty"`
  Generator         string          `json:"generator,omitempty"`
  Categories        string          `json:"categories,omitempty"`
  User              string          `json:"user,omitempty"`
  CreatedAt         time.Time       `json:"created_at,omitempty"`
  UpdatedAt         time.Time       `json:"update_at,omitempty"`
}
