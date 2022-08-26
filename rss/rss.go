package rss

import (
	// log "github.com/sirupsen/logrus"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"strings"

	"github.com/mrusme/journalist/crawler"
	"github.com/mrusme/journalist/ent"

	"github.com/araddon/dateparse"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

type Client struct {
  parser        *gofeed.Parser
  url           string
  username      string
  password      string
  Feed          *gofeed.Feed
  Items         *[]*gofeed.Item
  ItemsCrawled  []crawler.ItemCrawled
  exceptItemGUIDs []string
  UpdatedAt     time.Time
  logger        *zap.Logger
}

func NewClient(
  feedUrl string,
  username string,
  password string,
  crawl bool,
  exceptItemGUIDs []string,
  logger *zap.Logger,
) (*Client, []error) {
  client := new(Client)
  client.parser = gofeed.NewParser()
  client.url = feedUrl
  client.username = username
  client.password = password
  client.exceptItemGUIDs = exceptItemGUIDs
  client.logger = logger

  if errs := client.Sync(crawl); errs != nil {
    return nil, errs
  }

  return client, nil
}

func (c *Client) Sync(crawl bool) ([]error) {
  var errs []error

  c.logger.Debug(
    "Starting RSS Sync procedure",
    zap.Bool("crawl", crawl),
  )

  feedCrwl := crawler.New(c.logger)
  defer feedCrwl.Close()
  feedCrwl.SetLocation(c.url)
  feedCrwl.SetBasicAuth(c.username, c.password)
  feed, err := feedCrwl.ParseFeed()
  if err != nil {
    c.logger.Debug(
      "RSS Sync error occurred for feed crawling",
      zap.String("url", c.url),
      zap.Error(err),
    )
    errs = append(errs, err)
    return errs
  }

  c.Feed = feed
  c.Items = &feed.Items
  c.UpdatedAt = time.Now()

  if crawl == true {
    c.logger.Debug(
      "RSS Sync starting crawling procedure",
      zap.Int("exceptItemGUIDsLength", len(c.exceptItemGUIDs)),
    )
    crwl := crawler.New(c.logger)
    defer crwl.Close()

    for i := 0; i < len(c.Feed.Items); i++ {
      var foundException bool = false
      itemGUID := GenerateGUIDForItem(c.Feed.Items[i])
      for _, exceptItemGUID := range c.exceptItemGUIDs {
        if exceptItemGUID == itemGUID {
          c.logger.Debug(
            "Crawler found exception, breaking",
            zap.String("itemGUID", exceptItemGUID),
            zap.String("itemLink", c.Feed.Items[i].Link),
          )
          foundException = true
          break
        }
      }

      if foundException == true {
        continue
      }
      c.logger.Debug(
        "Crawler found no exception, continuing with item",
        zap.String("itemLink", c.Feed.Items[i].Link),
      )
      crwl.Reset()
      crwl.SetLocation(c.Feed.Items[i].Link)
      crwl.SetBasicAuth(c.username, c.password)
      itemCrawled, err := crwl.GetReadable()
      if err != nil {
        c.logger.Debug(
          "Crawler failed to GetReadable",
          zap.String("itemLink", c.Feed.Items[i].Link),
          zap.Error(err),
        )
        errs = append(errs, err)
        continue
      }

      c.ItemsCrawled = append(c.ItemsCrawled, itemCrawled)
    }
  }

  return errs
}

func (c* Client) SetFeed(
  feedLink string,
  username string,
  password string,
  dbFeedTmp *ent.FeedCreate,
) (*ent.FeedCreate) {
  // TODO: Get system timezone
  ltz, _ := time.LoadLocation("UTC")
  time.Local = ltz

  feedUpdated, err := dateparse.ParseLocal(c.Feed.Updated)
  if err != nil {
    feedUpdated = time.Now()
  }
  feedPublished, err := dateparse.ParseLocal(c.Feed.Published)
  if err != nil {
    feedPublished = time.Now()
  }

  dbFeedTmp = dbFeedTmp.
    SetURL(feedLink).
    SetUsername(username).
    SetPassword(password).
    SetFeedTitle(c.Feed.Title).
    SetFeedDescription(c.Feed.Description).
    SetFeedLink(c.Feed.Link).
    SetFeedFeedLink(c.Feed.FeedLink).
    SetFeedUpdated(feedUpdated).
    SetFeedPublished(feedPublished).
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

  // TODO: Get system timezone
  ltz, _ := time.LoadLocation("UTC")
  time.Local = ltz

  itemUpdated, err := dateparse.ParseLocal(item.Updated)
  if err != nil {
    itemUpdated = time.Now()
  }
  itemPublished, err := dateparse.ParseLocal(item.Published)
  if err != nil {
    itemPublished = time.Now()
  }

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
    SetItemUpdated(itemUpdated).
    SetItemPublished(itemPublished).
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

