package api

import (
  "net/http"
  "strconv"
)

func GetSinceIDFromReq(r *http.Request) (int64) {
  var sinceID int64

  _, hasSinceID := r.Form["since_id"]
  if hasSinceID == true {
    sinceID, _ = strconv.ParseInt(r.FormValue("since_id"), 10, 64)
  } else {
    sinceID = 0
  }

  return sinceID
}
