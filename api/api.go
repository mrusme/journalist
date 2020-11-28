package api

import (
  "net/http"
  "log"
  "io"
  "github.com/gorilla/mux"
)

func api(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if r.Method == http.MethodOptions {
      return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)

  io.WriteString(w, `{"alive": true}`)
}

func Server() {
  r := mux.NewRouter()
  r.HandleFunc("/", api)
  r.Use(mux.CORSMethodMiddleware(r))
  log.Fatal(http.ListenAndServe(":8000", r))
}
