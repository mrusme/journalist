package db

import (
  "time"
)

type Item struct {
  ID                uint            `json:"id,omitempty"`
  GUID              string          `json:"guid,omitempty"`
  Title             string          `json:"title,omitempty"`
  Description       string          `json:"description,omitempty"`
  Content           string          `json:"content,omitempty"`
  Link              string          `json:"link,omitempty"`
  Author            string          `json:"author,omitempty"`
  Image             string          `json:"image,omitempty"`
  Categories        string          `json:"categories,omitempty"`
  IsRead            bool            `json:"is_read,omitempty"`
  IsSaved           bool            `json:"is_saved,omitempty"`
  Feed              uint            `json:"feed,omitempty"`
  User              string          `json:"user,omitempty"`
  UpdatedAt         time.Time       `json:"updated_at,omitempty"`
  CreatedAt         time.Time       `json:"created_at,omitempty"`
}
