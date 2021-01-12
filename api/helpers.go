package api

import (
  "net/http"
  "strconv"
  log "github.com/sirupsen/logrus"
  "github.com/mrusme/journalist/db"
  "github.com/mrusme/journalist/rss"
  "github.com/mrusme/journalist/web"
)

func AddOrUpdateFeed(db *db.Database, feedLink string, feedGroup int64, feedUser string) (error) {
  var err error

  refreshedFeed, items, err := rss.LoadFeed(feedLink, feedGroup, feedUser)
  if err != nil {
    log.Error(err)
    return err
  }

  newItems, err := db.UpsertFeed(&refreshedFeed, items)
  if err != nil {
    log.Error(err)
    return err
  }

  log.Debug("Iterating through newly created items ...")
  for _, newItem := range newItems {
    log.Debug("New item with ID: ", newItem.ID)

    err = web.MakeItemReadable(newItem)
    if err != nil {
      log.Debug("Could not make " + newItem.Link + " readable!")
    }

    log.Debug("Updating item with readable content ...")
    err = db.UpdateItem(newItem)
    if err != nil {
      log.Error("Could not update item: ", newItem.ID)
      log.Error(err)
    }
  }

  return err
}

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
