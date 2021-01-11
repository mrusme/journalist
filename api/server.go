package api

import (
  "os"
  "strconv"
  "net/http"
  "time"
  log "github.com/sirupsen/logrus"
  "github.com/gorilla/mux"
  "github.com/mrusme/journalist/db"
  "github.com/mrusme/journalist/rss"
)

var database *db.Database

func Server(db *db.Database) {
  database = db

  portStr, ok := os.LookupEnv("JOURNALIST_SERVER_PORT")
  if ok == false {
    portStr = "8000"
  }

  refreshStr, ok := os.LookupEnv("JOURNALIST_SERVER_REFRESH")
  if ok == false {
    refreshStr = "0"
  }

  refresh, parseerr := strconv.ParseInt(refreshStr, 10, 64)
  if parseerr != nil {
    log.Fatal(parseerr)
  }

  if refresh > 0 {
    go refreshLoop(db, refresh)
  }

  r := mux.NewRouter()
  r.Use(mux.CORSMethodMiddleware(r))

  feverAPIRouter := r.PathPrefix("/fever").Subrouter()
  feverAPI(feverAPIRouter)

  greaderAPIRouter := r.PathPrefix("/greader").Subrouter()
  greaderAPI(greaderAPIRouter)

  log.Info("Starting server on port " + portStr + " ...")

  server := &http.Server{
    Addr:         "0.0.0.0:" + portStr,
    WriteTimeout: time.Second * 60,
    ReadTimeout:  time.Second * 60,
    IdleTimeout:  time.Second * 80,
    Handler: r,
  }
  log.Fatal(server.ListenAndServe())
}

func refreshLoop(db *db.Database, interval int64) {
  intervalDuration := time.Second * time.Duration(interval)

  for {
    refresh(db)
    time.Sleep(intervalDuration)
  }
}

func refresh(db *db.Database) {

  log.Debug("Refreshing feeds ...")
  feeds, err := db.ListFeeds()
  if err != nil {
    log.Error(err)
    return
  }

  for _, feed := range feeds {
    log.Debug("Refreshing ", feed.FeedLink, " ...")

    refreshedFeed, items, feederr := rss.LoadFeed(feed.FeedLink, feed.Group, feed.User)
    if feederr != nil {
      log.Error(feederr)
      return
    }

    _, upserterr := database.UpsertFeed(&refreshedFeed, &items)
    if upserterr != nil {
      log.Error(upserterr)
      return
    }
  }

  log.Debug("Refresh completed")
}
