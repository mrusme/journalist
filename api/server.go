package api

import (
  "net/http"
  "time"
  log "github.com/sirupsen/logrus"
  "github.com/gorilla/mux"
  "github.com/mrusme/journalist/db"
  "github.com/mrusme/journalist/common"
)

var database *db.Database

func Server(db *db.Database) {
  database = db

  bindIPStr := common.LookupStrEnv("JOURNALIST_SERVER_BINDIP", "0.0.0.0")
  portStr := common.LookupStrEnv("JOURNALIST_SERVER_BINDIP", "8000")
  refresh := common.LookupIntEnv("JOURNALIST_SERVER_REFRESH", 0)

  if refresh > 0 {
    go refreshLoop(db, refresh)
  }

  r := mux.NewRouter()
  r.Use(mux.CORSMethodMiddleware(r))

  if common.LookupBooleanEnv("JOURNALIST_SERVER_API_FEVER", true) == true {
    log.Info("Enabling Fever API ...")
    feverAPIRouter := r.PathPrefix("/fever").Subrouter()
    feverAPI(feverAPIRouter)
  }

  if common.LookupBooleanEnv("JOURNALIST_SERVER_API_GREADER", false) == true {
    log.Info("Enabling Google Reader API ...")
    greaderAPIRouter := r.PathPrefix("/greader").Subrouter()
    greaderAPI(greaderAPIRouter)
  }

  log.Info("Starting server on " + bindIPStr + ":" + portStr + " ...")

  handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    log.Println(req.Method + " " + req.URL.String())
    r.ServeHTTP(w, req)
  })

  server := &http.Server{
    Addr:         bindIPStr + ":" + portStr,
    WriteTimeout: time.Second * 60,
    ReadTimeout:  time.Second * 60,
    IdleTimeout:  time.Second * 80,
    Handler: handler,
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

    AddOrUpdateFeed(db, feed.FeedLink, feed.Group, feed.User)
  }

  log.Debug("Refresh completed")
}
