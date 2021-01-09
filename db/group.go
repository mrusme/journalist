package db

import (
  "time"
)

type Group struct {
  ID                uint            `db:"id",json:"id,omitempty"`
  Title             string          `db:"title",json:"title,omitempty"`
  TitleUnix         string          `db:"title_unix",json:"title_unix,omitempty"`
  User              string          `db:"user",json:"user,omitempty"`
  CreatedAt         time.Time       `db:"created_at",json:"created_at,omitempty"`
  UpdatedAt         time.Time       `db:"updated_at",json:"updated_at,omitempty"`
}
