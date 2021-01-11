package api

import (
  "net/http"
  // log "github.com/sirupsen/logrus"
  // "encoding/json"
  "github.com/gorilla/mux"
  // "github.com/mrusme/journalist/db"
)

func greaderAPIHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if r.Method == http.MethodOptions {
      return
  }

  // TODO: Implement Google Reader API
  // See:
  // - http://code.google.com/p/pyrfeed/wiki/GoogleReaderAPI
  // - https://blog.martindoms.com/2009/10/16/using-the-google-reader-api-part-2
  // - https://ranchero.com/downloads/GoogleReaderAPI-2009.pdf
  // - https://github.com/theoldreader/api
  // - https://github.com/devongovett/reader
  // - https://github.com/FreshRSS/FreshRSS/blob/master/p/api/greader.php

  w.WriteHeader(http.StatusNoContent)
  return
}

func greaderAPI(r *mux.Router) {
  r.HandleFunc("/", greaderAPIHandler)
}
