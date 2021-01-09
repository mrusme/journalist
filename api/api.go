package api

import (
  "net/http"
  log "github.com/sirupsen/logrus"
  "errors"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/mrusme/journalist/db"
)

var database *db.Database

type ApiFeedsGroup struct {
  FeedIDs             string          `json:"feed_ids,omitempty"`
  GroupID             uint            `json:"group_id,omitempty"`
}

type ApiGroup struct {
  ID                  int64           `json:"id,omitempty"`
  Title               string          `json:"title,omitempty"`
}

type ApiFeed struct {
  ID                  int64           `json:"id,omitempty"`
  Title               string          `json:"title,omitempty"`
  SiteURL             string          `json:"site_url,omitempty"`
  URL                 string          `json:"url,omitempty"`
  LastUpdatedOnTime   int             `json:"last_updated_on_time,omitempty"`
  IsSpark             bool            `json:"is_spark,omitempty"`
}

type ApiFavicon struct {
  ID                  int64           `json:"id,omitempty"`
  Data                string          `json:"data,omitempty"`
}

type ApiItem struct {
  ID                  int64           `json:"id,omitempty"`
  FeedID              int64           `json:"feed_id,omitempty"`
  Title               string          `json:"title,omitempty"`
  URL                 string          `json:"url,omitempty"`
  Author              string          `json:"author,omitempty"`
  HTML                string          `json:"html,omitempty"`
  CreatedOnTime       int             `json:"created_on_time,omitempty"`
  IsRead              bool            `json:"is_read,omitempty"`
  IsSaved             bool            `json:"is_saved,omitempty"`
}

type ApiResponse struct {
  ApiVersion          string          `json:"api_version,omitempty"`
  Auth                int             `json:"auth,omitempty"`
  FeedsGroups         []ApiFeedsGroup `json:"feeds_groups,omitempty"`
  Groups              []ApiGroup      `json:"groups,omitempty"`
  Feeds               []ApiFeed       `json:"feeds,omitempty"`
  Favicons            []ApiFavicon    `json:"favicons,omitempty"`
  Items               []ApiItem       `json:"items,omitempty"`
  TotalItems          int             `json:"total_items,omitempty"`
  UnreadItemIDs       string          `json:"unread_item_ids,omitempty"`
  SavedItemIDs        string          `json:"saved_item_ids,omitempty"`
  LastRefreshedOnTime int             `json:"last_refreshed_on_time,omitempty"`
}

func (apiResponse *ApiResponse) processGroups(r *http.Request, user string) (bool, error) {
  _, hasGroups := r.Form["groups"]
  if hasGroups == true {
    groups, err := database.ListGroupsByUser(user)
    if err != nil {
      log.Error(err)
    }

    for _, group := range groups {
      apiResponse.Groups = append(apiResponse.Groups,
        ApiGroup{
          ID: group.ID,
          Title: group.Title,
        })
    }
    return true, err
  }

  return false, nil
}

func (apiResponse *ApiResponse) processFeeds(r *http.Request, user string) (bool, error) {
  _, hasFeeds := r.Form["feeds"]
  if hasFeeds == true {
    feeds, err := database.ListFeedsByUser(user)
    if err != nil {
      log.Error(err)
    }

    for _, feed := range feeds {
      apiResponse.Feeds = append(apiResponse.Feeds,
        ApiFeed{
          ID: feed.ID,
          Title: feed.Title,
          SiteURL: feed.Link,
          URL: feed.FeedLink,
          LastUpdatedOnTime: int(feed.UpdatedAt.Unix()),
          IsSpark: false,
        })
    }
    return true, err
  }

  return false, nil
}

func (apiResponse *ApiResponse) processItems(r *http.Request, user string) (bool, error) {
  _, hasItems := r.Form["items"]
  if hasItems == true {
    // TODO
    return true, nil
  }

  return false, nil
}

func (apiResponse *ApiResponse) processUnreadItemIDs(r *http.Request, user string) (bool, error) {
  _, hasUnreadItemIDs := r.Form["unread_item_ids"]
  if hasUnreadItemIDs == true {
    apiResponse.UnreadItemIDs = "1,2,3"
    return true, nil
  }

  return false, nil
}

func (apiResponse *ApiResponse) processSavedItemIDs(r *http.Request, user string) (bool, error) {
  _, hasSavedItemIDs := r.Form["saved_item_ids"]
  if hasSavedItemIDs == true {
    apiResponse.SavedItemIDs = "1,2,3"
    return true, nil
  }

  return false, nil
}

func (apiResponse *ApiResponse) processMark(r *http.Request, user string) (bool, error) {
  _, hasMark := r.Form["mark"]
  if hasMark == true {
    mark := r.FormValue("mark")
    if mark != "item" && mark != "feed" && mark != "group" {
      return false, errors.New("`mark` parameter must be either `item`, `feed` or `group`")
    }

    _, hasAs := r.Form["as"]
    if hasAs == false {
      return false, errors.New("`as` parameter missing")
    }

    as := r.FormValue("as")
    if as != "read" && as != "saved" && as != "unsaved" {
      return false, errors.New("`as` parameter must be either `read`, `saved` or `unsaved`")
    }

    _, hasId := r.Form["id"]
    if hasId == false {
      return false, errors.New("`id` parameter missing")
    }

    _, hasBefore := r.Form["before"]
    if (mark == "feed" || mark == "group") && hasBefore == false {
      return false, errors.New("`before` parameter missing")
    }

    // TODO

    return true, nil
  }

  return false, nil
}

func api(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if r.Method == http.MethodOptions {
      return
  }

  r.ParseMultipartForm(8192)
  log.Printf("%+v", r.Form)
  log.Printf("%+v", r.PostForm)

  user := r.PostFormValue("api_key")
  if user == "" {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  _, hasApi := r.Form["api"]
  if hasApi == false {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  apiResponse := ApiResponse{}
  apiResponse.ApiVersion = "4"
  apiResponse.Auth = 1

  // var processed bool
  var processError error

  // Groups
  _, processError = apiResponse.processGroups(r, user)
  if processError != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // Feeds
  _, processError = apiResponse.processFeeds(r, user)
  if processError != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // Items
  _, processError = apiResponse.processItems(r, user)
  if processError != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // Unread item IDs
  _, processError = apiResponse.processUnreadItemIDs(r, user)
  if processError != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // Saved item IDs
  _, processError = apiResponse.processSavedItemIDs(r, user)
  if processError != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  // Mark ... as ...
  _, processError = apiResponse.processMark(r, user)
  if processError != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)

  json.NewEncoder(w).Encode(apiResponse)
}

func Server(db *db.Database) {
  database = db
  r := mux.NewRouter()
  r.HandleFunc("/", api)
  r.Use(mux.CORSMethodMiddleware(r))
  log.Fatal(http.ListenAndServe(":8000", r))
}
