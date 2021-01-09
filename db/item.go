package db

import (
  "time"
)

type Item struct {
  ID                uint            `db:"id",json:"id,omitempty"`
  GUID              string          `db:"guid",json:"guid,omitempty"`
  Title             string          `db:"title",json:"title,omitempty"`
  Description       string          `db:"description",json:"description,omitempty"`
  Content           string          `db:"content",json:"content,omitempty"`
  Link              string          `db:"link",json:"link,omitempty"`
  Author            string          `db:"author",json:"author,omitempty"`
  Image             string          `db:"image",json:"image,omitempty"`
  Categories        string          `db:"categories",json:"categories,omitempty"`
  IsRead            bool            `db:"is_read",json:"is_read,omitempty"`
  IsSaved           bool            `db:"is_saved",json:"is_saved,omitempty"`
  Feed              uint            `db:"feed",json:"feed,omitempty"`
  User              string          `db:"user",json:"user,omitempty"`
  CreatedAt         time.Time       `db:"created_at",json:"created_at,omitempty"`
  UpdatedAt         time.Time       `db:"updated_at",json:"updated_at,omitempty"`
}
