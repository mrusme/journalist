package web

import (
  log "github.com/sirupsen/logrus"
  "net/http"
  readability "github.com/go-shiori/go-readability"
)

func LoadReadablePage(pageUrl string) (readability.Article, error) {
  log.Debug("Making " + pageUrl + " readable ...")

  resp, err := http.Get(pageUrl)
  if err != nil {
    log.Debug("Failed to get page: " + pageUrl)
    return readability.Article{}, err
  }
  defer resp.Body.Close()

  article, err := readability.FromReader(resp.Body, pageUrl)
  if err != nil {
    return readability.Article{}, err
  }

  return article, nil
}
