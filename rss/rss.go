package rss

import (
	// log "github.com/sirupsen/logrus"
	"github.com/google/uuid"
	"crypto/sha256"
  "encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"strings"

	"github.com/mrusme/journalist/crawler"
	"github.com/mrusme/journalist/ent"

	"github.com/mmcdole/gofeed"
  "github.com/microcosm-cc/bluemonday"
)

type Client struct {
  parser        *gofeed.Parser
  url           string
  username      string
  password      string
  Feed          *gofeed.Feed
  Items         *[]*gofeed.Item
  ItemsCrawled  []crawler.ItemCrawled
  UpdatedAt     time.Time
}

func NewClient(
  feedUrl string,
  username string,
  password string,
  crawl bool,
) (*Client, error) {
  client := new(Client)
  client.parser = gofeed.NewParser()
  client.url = feedUrl
  client.username = username
  client.password = password

  if err := client.Sync(crawl); err != nil {
    return nil, err
  }

  return client, nil
}

func (c *Client) Sync(crawl bool) (error) {
  var errs []error

  feedCrwl := crawler.New()
  defer feedCrwl.Close()
  feedCrwl.SetLocation(c.url)
  feedCrwl.SetBasicAuth(c.username, c.password)
  feed, err := feedCrwl.ParseFeed()
  if err != nil {
    return err
  }

  c.Feed = feed
  c.Items = &feed.Items
  c.UpdatedAt = time.Now()

  if crawl == true {
    crwl := crawler.New()
    defer crwl.Close()

    for i := 0; i < len(c.Feed.Items); i++ {
      crwl.Reset()
      crwl.SetLocation(c.Feed.Items[i].Link)
      crwl.SetBasicAuth(c.username, c.password)
      itemCrawled, err := crwl.GetReadable()
      if err != nil {
        errs = append(errs, err)
        continue
      }

      c.ItemsCrawled = append(c.ItemsCrawled, itemCrawled)
    }
  }

  return nil
}

func (c* Client) SetFeed(
  feedLink string,
  username string,
  password string,
  dbFeedTmp *ent.FeedCreate,
) (*ent.FeedCreate) {
  dbFeedTmp = dbFeedTmp.
    SetURL(feedLink).
    SetUsername(username).
    SetPassword(password).
    SetFeedTitle(c.Feed.Title).
    SetFeedDescription(c.Feed.Description).
    SetFeedLink(c.Feed.Link).
    SetFeedFeedLink(c.Feed.FeedLink).
    SetFeedUpdated(c.Feed.Updated).
    SetFeedPublished(c.Feed.Published).
    SetFeedLanguage(c.Feed.Language).
    SetFeedCopyright(c.Feed.Copyright).
    SetFeedGenerator(c.Feed.Generator).
    SetFeedCopyright(c.Feed.Copyright).
    SetFeedCategories(strings.Join(c.Feed.Categories, ", "))

  if c.Feed.Author != nil {
    dbFeedTmp = dbFeedTmp.
      SetFeedAuthorName(c.Feed.Author.Name).
      SetFeedAuthorEmail(c.Feed.Author.Email)
  }
  if c.Feed.Image != nil {
    dbFeedTmp = dbFeedTmp.
      SetFeedImageTitle(c.Feed.Image.Title).
      SetFeedImageURL(c.Feed.Image.URL)
  }

  return dbFeedTmp
}

func (c* Client) SetItem(
  feedID uuid.UUID,
  idx int,
  dbItemTemp *ent.ItemCreate,
) (*ent.ItemCreate) {
  var crawled crawler.ItemCrawled
  if len(c.ItemsCrawled) > idx {
    crawled = c.ItemsCrawled[idx]
  }

  item := c.Feed.Items[idx]

  var enclosureJson string = ""
  if item.Enclosures != nil {
    jsonbytes, err := json.Marshal(item.Enclosures)
    if err == nil {
      enclosureJson = string(jsonbytes)
    }
  }

  itemDescription := bluemonday.
    StrictPolicy().
    Sanitize(item.Description)

  dbItemTemp = dbItemTemp.
    SetFeedID(feedID).
    SetItemGUID(GenerateGUIDForItem(item)).
    SetItemTitle(item.Title).
    SetItemDescription(itemDescription).
    SetItemContent(item.Content).
    SetItemLink(item.Link).
    SetItemUpdated(item.Updated).
    SetItemPublished(item.Published).
    SetItemCategories(strings.Join(item.Categories, ",")).
    SetItemEnclosures(enclosureJson).

    SetCrawlerTitle(crawled.Title).
    SetCrawlerAuthor(crawled.Author).
    SetCrawlerExcerpt(crawled.Excerpt).
    SetCrawlerSiteName(crawled.SiteName).
    SetCrawlerImage(crawled.Image).
    SetCrawlerContentHTML(crawled.ContentHtml).
    SetCrawlerContentText(crawled.ContentText)

  if item.Author != nil {
    dbItemTemp = dbItemTemp.
      SetItemAuthorName(item.Author.Name).
      SetItemAuthorEmail(item.Author.Email)
  }

  if item.Image != nil {
    dbItemTemp = dbItemTemp.
      SetItemImageTitle(item.Image.Title).
      SetItemImageURL(item.Image.URL)
  }

  return dbItemTemp
}

func GenerateGUID(from string) (string) {
  h := sha256.New()
  h.Write([]byte(from))
  return hex.EncodeToString(
    h.Sum(nil),
  )
}

func GenerateGUIDForItem(item *gofeed.Item) (string) {
  return GenerateGUID(
    fmt.Sprintf("%s%s", item.Link, item.Published),
  )
}

