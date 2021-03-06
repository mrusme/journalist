package api

import (
  "net/http"
  log "github.com/sirupsen/logrus"
  "strings"
  "strconv"
  "errors"
  "time"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/mrusme/journalist/db"
)

type FeverAPIFeedsGroup struct {
  FeedIDs             string          `json:"feed_ids,omitempty"`
  GroupID             int64           `json:"group_id,omitempty"`
}

type FeverAPIGroup struct {
  ID                  int64           `json:"id,omitempty"`
  Title               string          `json:"title,omitempty"`
}

type FeverAPIFeed struct {
  ID                  int64           `json:"id,omitempty"`
  Title               string          `json:"title,omitempty"`
  SiteURL             string          `json:"site_url,omitempty"`
  URL                 string          `json:"url,omitempty"`
  LastUpdatedOnTime   int             `json:"last_updated_on_time,omitempty"`
  IsSpark             bool            `json:"is_spark,omitempty"`
}

type FeverAPIFavicon struct {
  ID                  int64           `json:"id,omitempty"`
  Data                string          `json:"data,omitempty"`
}

type FeverAPIItem struct {
  ID                  int64           `json:"id,omitempty"`
  FeedID              int64           `json:"feed_id,omitempty"`
  Title               string          `json:"title,omitempty"`
  URL                 string          `json:"url,omitempty"`
  Author              string          `json:"author,omitempty"`
  HTML                string          `json:"html,omitempty"`
  CreatedOnTime       int             `json:"created_on_time,omitempty"`
  IsRead              int             `json:"is_read"`
  IsSaved             int             `json:"is_saved"`
}

type FeverAPIResponse struct {
  ApiVersion          string          `json:"api_version,omitempty"`
  Auth                int             `json:"auth,omitempty"`
  FeedsGroups         []FeverAPIFeedsGroup `json:"feeds_groups,omitempty"`
  Groups              []FeverAPIGroup      `json:"groups,omitempty"`
  Feeds               []FeverAPIFeed       `json:"feeds,omitempty"`
  Favicons            []FeverAPIFavicon    `json:"favicons,omitempty"`
  Items               []FeverAPIItem       `json:"items,omitempty"`
  TotalItems          int             `json:"total_items,omitempty"`
  UnreadItemIDs       string          `json:"unread_item_ids,omitempty"`
  SavedItemIDs        string          `json:"saved_item_ids,omitempty"`
  LastRefreshedOnTime int             `json:"last_refreshed_on_time,omitempty"`
}

func (apiResponse *FeverAPIResponse) processGroups(r *http.Request, user string) (bool, error) {
  _, hasGroups := r.Form["groups"]
  if hasGroups == true {
    groups, err := database.ListGroupsByUser(user)
    if err != nil {
      log.Error(err)
    }

    for _, group := range groups {
      apiResponse.Groups = append(apiResponse.Groups,
        FeverAPIGroup{
          ID: group.ID,
          Title: group.Title,
        })

      feeds, err := database.ListFeedsByGroupAndUser(group.ID, user)
      if err != nil {
        return false, err
      }

      var feedIDsStr []string
      for _, feed := range feeds {
        feedIDsStr = append(feedIDsStr, strconv.FormatInt(feed.ID, 10))
      }

      apiResponse.FeedsGroups = append(apiResponse.FeedsGroups,
        FeverAPIFeedsGroup{
          GroupID: group.ID,
          FeedIDs: strings.Join(feedIDsStr, ","),
        })
    }

    return true, err
  }

  return false, nil
}

func (apiResponse *FeverAPIResponse) processFeeds(r *http.Request, user string) (bool, error) {
  _, hasFeeds := r.Form["feeds"]
  if hasFeeds == true {
    feeds, err := database.ListFeedsByUser(user)
    if err != nil {
      log.Error(err)
    }

    for _, feed := range feeds {
      log.Debug(feed.UpdatedAt.Unix())
      apiResponse.Feeds = append(apiResponse.Feeds,
        FeverAPIFeed{
          ID: feed.ID,
          Title: feed.Title,
          SiteURL: feed.Link,
          URL: feed.FeedLink,
          LastUpdatedOnTime: int(feed.UpdatedAt.Unix()),
          IsSpark: false,
        })
    }
    log.Debug("Returning feeds ...")
    return true, err
  }

  return false, nil
}

func (apiResponse *FeverAPIResponse) processItems(r *http.Request, user string) (bool, error) {
  _, hasItems := r.Form["items"]
  if hasItems == true {
    var items []db.Item
    var err error

    _, hasWithIDs := r.Form["with_ids"]
    if hasWithIDs == true {
      withIDsStr := strings.Split(r.FormValue("with_ids"), ",")
      var withIDs []int64
      for _, withIDStr := range withIDsStr {
        i, _ := strconv.ParseInt(withIDStr, 10, 64)
        withIDs = append(withIDs, i)
      }

      items, err = database.ListItemsByIDsAndUser(withIDs, user)
      if err != nil {
        log.Error(err)
      }
    } else {
      sinceID := GetSinceIDFromReq(r)
      items, err = database.ListItemsByUser(user, sinceID)
      if err != nil {
        log.Error(err)
      }
    }

    for _, item := range items {
      isRead := 0
      if item.IsRead == true {
        isRead = 1
      }

      isSaved := 0
      if item.IsSaved == true {
        isSaved = 1
      }

      var html *string
      ReadableContentLen := len(item.ReadableContent)
      DescriptionLen := len(item.Description)
      if item.ReadableContent != "" && ReadableContentLen >= DescriptionLen {
        html = &item.ReadableContent
      } else {
        html = &item.Description
      }

      apiResponse.Items = append(apiResponse.Items,
        FeverAPIItem{
          ID: item.ID,
          FeedID: item.Feed,
          Title: item.Title,
          URL: item.Link,
          Author: item.Author,
          HTML: *html,
          CreatedOnTime: int(item.CreatedAt.Unix()),
          IsRead: isRead,
          IsSaved: isSaved,
        })
    }
    log.Debug("Returning items ...")
    return true, nil
  }

  return false, nil
}

func (apiResponse *FeverAPIResponse) processUnreadItemIDs(r *http.Request, user string) (bool, error) {
  _, hasUnreadItemIDs := r.Form["unread_item_ids"]
  if hasUnreadItemIDs == true {
    sinceID := GetSinceIDFromReq(r)
    items, err := database.ListUnreadItemsByUser(user, sinceID)
    if err != nil {
      log.Error(err)
    }

    var itemIDs []string
    for _, item := range items {
      itemIDs = append(itemIDs, strconv.FormatInt(item.ID, 10))
    }

    apiResponse.UnreadItemIDs = strings.Join(itemIDs, ",")
    log.Debug("Returning unread item IDs ...")
    return true, nil
  }

  return false, nil
}

func (apiResponse *FeverAPIResponse) processSavedItemIDs(r *http.Request, user string) (bool, error) {
  _, hasSavedItemIDs := r.Form["saved_item_ids"]
  if hasSavedItemIDs == true {
    apiResponse.SavedItemIDs = ""
    return true, nil
  }

  return false, nil
}

func (apiResponse *FeverAPIResponse) processMark(r *http.Request, user string) (bool, error) {
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
    if as != "read" && as != "unread" && as != "saved" && as != "unsaved" {
      return false, errors.New("`as` parameter must be either `read`, `unread`, `saved` or `unsaved`")
    }

    _, hasId := r.Form["id"]
    if hasId == false {
      return false, errors.New("`id` parameter missing")
    }

    idStr := r.FormValue("id")
    id, _ := strconv.ParseInt(idStr, 10, 64)

    _, hasBefore := r.Form["before"]
    if (mark == "feed" || mark == "group") && hasBefore == false {
      return false, errors.New("`before` parameter missing")
    }

    var before time.Time
    if hasBefore == true {
      beforeStr := r.FormValue("before")
      beforeInt, _ := strconv.ParseInt(beforeStr, 10, 64)
      before = time.Unix(beforeInt, 0)
    }

    var err error
    // Mark item with ID x as read
    if mark == "item" && hasId == true {
      switch(as) {
      case "read":
        err = database.UpdateItemByIDAsRead(id, true, user)
      case "unread":
        err = database.UpdateItemByIDAsRead(id, false, user)
      case "saved":
        err = database.UpdateItemByIDAsSaved(id, true, user)
      case "unsaved":
        err = database.UpdateItemByIDAsSaved(id, false, user)
      }
    } else if mark == "feed" && hasBefore == true {
      switch(as) {
      case "read":
        err = database.UpdateItemsByFeedAndBeforeAsRead(id, before, true, user)
      }
    } else if mark == "group" && hasBefore == true {
      switch(as) {
      case "read":
        if id > 0 {
          err = database.UpdateItemsByGroupAsRead(id, true, user)
        } else {
          err = database.UpdateItemsByBeforeAsRead(before, true, user)
        }
      }
    }

    if err != nil {
      log.Debug(err)
      return false, err
    }

    return true, nil
  }

  return false, nil
}

func feverAPIHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  if r.Method == http.MethodOptions {
      return
  }

  r.ParseMultipartForm(8192)
  log.Printf("%+v", r.Form)
  log.Printf("%+v", r.PostForm)

  _, hasRefresh := r.Form["refresh"]
  if hasRefresh == true {
    refresh(database)
    w.WriteHeader(http.StatusNoContent)
    return
  }

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

  apiResponse := FeverAPIResponse{}
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

func feverAPI(r *mux.Router) {
  r.HandleFunc("/", feverAPIHandler)
}
