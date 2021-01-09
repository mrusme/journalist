package db

import (
  "time"
)

type Feed struct {
  ID                int64           `db:"id",json:"id,omitempty"`
  Title             string          `db:"title",json:"title,omitempty"`
  Description       string          `db:"description",json:"description,omitempty"`
  Link              string          `db:"link",json:"link,omitempty"`
  FeedLink          string          `db:"feed_link",json:"feed_link,omitempty"`
  Author            string          `db:"author",json:"author,omitempty"`
  Language          string          `db:"language",json:"language,omitempty"`
  Image             string          `db:"image",json:"image,omitempty"`
  Copyright         string          `db:"copyright",json:"copyright,omitempty"`
  Generator         string          `db:"generator",json:"generator,omitempty"`
  Categories        string          `db:"categories",json:"categories,omitempty"`
  Group             int64           `db:"group",json:"group,omitempty"`
  User              string          `db:"user",json:"user,omitempty"`
  CreatedAt         time.Time       `db:"created_at",json:"created_at,omitempty"`
  UpdatedAt         time.Time       `db:"updated_at",json:"update_at,omitempty"`
}
