package db

import (
)

type Group struct {
  ID                uint            `json:"id,omitempty"`
  Title             string          `json:"title,omitempty"`
  User              string          `json:"user,omitempty"`
  CreatedAt         int64           `json:"created_at,omitempty"`
  UpdatedAt         int64           `json:"update_at,omitempty"`
}
