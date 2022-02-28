package web

import (
  log "github.com/sirupsen/logrus"
  "net/http"
  "net/url"
  readability "github.com/go-shiori/go-readability"
  "github.com/mrusme/journalist/db"
)

func MakeItemReadable(item *db.Item) (error) {
  pageUrl := item.Link
  log.Debug("Making " + pageUrl + " readable ...")

  resp, err := http.Get(pageUrl)
  if err != nil {
    log.Debug("Failed to get page: " + pageUrl)
    return err
  }
  defer resp.Body.Close()

  urlUrl, err := url.Parse(pageUrl)
  if err != nil {
    return err
  }

  article, err := readability.FromReader(resp.Body, urlUrl)
  if err != nil {
    return err
  }

  item.AssignReadableFromArticle(&article)

  return nil
}
