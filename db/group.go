package db

import (
)

type Group struct {
  ID                string          `json:"-"`
  IncID             int             `json:"inc_id,omitempty"`
  Title             string          `json:"title,omitempty"`
}

func (group *Group) SetIDFromDatabaseKey(key string) (error) {
  var err error
  group.ID, err = GetIDFromDatabaseKey(key)
  return err
}
