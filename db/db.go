package db

import (
  "strings"
  "regexp"
  "errors"
  log "github.com/sirupsen/logrus"

  _ "database/sql"
  "github.com/jmoiron/sqlx"
  _ "github.com/jackc/pgx/v4/stdlib"
  "github.com/mrusme/journalist/common"
)

var schema = `
CREATE TABLE IF NOT EXISTS groups (
    "id" SERIAL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "title_unix" TEXT NOT NULL,
    "user" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS "groups_title_unix" ON "groups"("title_unix","user");

CREATE TABLE IF NOT EXISTS feeds (
    "id" SERIAL PRIMARY KEY,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "link" TEXT NOT NULL,
    "feed_link" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "language" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "copyright" TEXT NOT NULL,
    "generator" TEXT NOT NULL,
    "categories" TEXT NOT NULL,
    "group" INT NOT NULL,
    "user" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT fk_group FOREIGN KEY("group") REFERENCES groups("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "feeds_feed_link" ON "feeds"("feed_link","user");

CREATE TABLE IF NOT EXISTS items (
    "id" SERIAL PRIMARY KEY,
    "guid" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "link" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "categories" TEXT NOT NULL,
    "readable_title" TEXT NOT NULL,
    "readable_author" TEXT NOT NULL,
    "readable_excerpt" TEXT NOT NULL,
    "readable_site_name" TEXT NOT NULL,
    "readable_image" TEXT NOT NULL,
    "readable_content" TEXT NOT NULL,
    "readable_text" TEXT NOT NULL,
    "is_read" BOOL NOT NULL,
    "is_saved" BOOL NOT NULL,
    "feed" INT NOT NULL,
    "user" TEXT NOT NULL,
    "created_at" TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP NOT NULL,
    CONSTRAINT fk_feed FOREIGN KEY("feed") REFERENCES feeds("id")
);
CREATE UNIQUE INDEX IF NOT EXISTS "items_guid" ON "items"("guid","user");
`

type Database struct {
  DB *sqlx.DB
}

func InitDatabase() (*Database, error) {
  dbconnection := common.LookupStrEnv("JOURNALIST_DB", "")
  if dbconnection == "" {
    return nil, errors.New("please `export JOURNALIST_DB` with the database connection string, e.g. 'postgres://user:secret@localhost:5432/journalist?sslmode=disable'")
  }

  log.Debug("Opening database connection ...")
  db, err := sqlx.Open("pgx", dbconnection)
  if err != nil {
    return nil, err
  }

  log.Debug("Connecting to database ...")
  err = db.Ping()
  if err != nil {
    db.Close()
    return nil, err
  }

  log.Debug("Creating database schema ...")
  db.MustExec(schema)

  database := Database{db}
  return &database, nil
}

func GetUnixName(name string) string {
  reg, regerr := regexp.Compile("[^a-zA-Z0-9]+")
  if regerr != nil {
      return ""
  }

  id := strings.ToLower(reg.ReplaceAllString(name, ""))

  return id
}
