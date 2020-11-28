package api

import (
  "net/http"
  "log"
  "io"
  "github.com/gorilla/mux"
)

type ApiFeedsGroup struct {
  FeedIDs        string          `json:"feed_ids,omitempty"`
  GroupID        int             `json:"group_id,omitempty"`
}

type ApiGroup struct {
  ID             int             `json:"id,omitempty"`
  Title          string          `json:"title,omitempty"`
}

type ApiFeed struct {
  ID                int             `json:"id,omitempty"`
  Title             string          `json:"title,omitempty"`
  SiteURL           string          `json:"site_url,omitempty"`
  URL               string          `json:"url,omitempty"`
  LastUpdatedOnTime int             `json:"last_updated_on_time,omitempty"`
  IsSpark           bool            `json:"is_spark,omitempty"`
}

type ApiFavicon struct {
  ID             int             `json:"id,omitempty"`
  Data           string          `json:"data,omitempty"`
}

type ApiItem struct {
  ID                int             `json:"id,omitempty"`
  FeedID            int             `json:"feed_id,omitempty"`
  Title             string          `json:"title,omitempty"`
  URL               string          `json:"url,omitempty"`
  Author            string          `json:"author,omitempty"`
  HTML              string          `json:"html,omitempty"`
  CreatedOnTime     int             `json:"created_on_time,omitempty"`
  IsRead            bool            `json:"is_read,omitempty"`
  IsSaved           bool            `json:"is_saved,omitempty"`
}

type ApiResponse struct {
  ApiVersion     string          `json:"api_version,omitempty"`
  Auth           int             `json:"auth,omitempty"`
  FeedsGroups    []ApiFeedsGroup `json:"feeds_groups,omitempty"`
  Groups         []ApiGroup      `json:"groups,omitempty"`
  Feeds          []ApiFeed       `json:"feeds,omitempty"`
  Favicons       []ApiFavicon    `json:"favicons,omitempty"`
  Items          []ApiItem       `json:"items,omitempty"`
}

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
