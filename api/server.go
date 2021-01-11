package api

import (
  "net/http"
  log "github.com/sirupsen/logrus"
  "github.com/gorilla/mux"
  "github.com/mrusme/journalist/db"
)

var database *db.Database

func Server(db *db.Database) {
  database = db

  go refreshLoop(db)

  r := mux.NewRouter()
  r.HandleFunc("/fever/", feverAPI)
  r.Use(mux.CORSMethodMiddleware(r))
  log.Fatal(http.ListenAndServe(":8000", r))
}
