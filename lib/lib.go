package lib

import (
  "go.uber.org/zap"

  "github.com/mrusme/journalist/ent"
)

type JournalistContext struct {
  Config                *Config
  EntClient             *ent.Client
  Logger                *zap.Logger
}

