package journalistd

import (
	"context"

	"github.com/google/uuid"

	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/feed"

	"github.com/mrusme/journalist/rss"
)

type Journalistd struct {
  entClient             *ent.Client
}

func New(entClient *ent.Client) (*Journalistd) {
  jd := new(Journalistd)
  jd.entClient = entClient
  return jd
}

func (jd *Journalistd) RefreshAll() ([]error) {
  var errs []error

  dbFeeds, err := jd.entClient.Feed.
    Query().
    All(context.Background())
  if err != nil {
    errs = append(errs, err)
    return errs
  }

  var feedIds []uuid.UUID = make([]uuid.UUID, len(dbFeeds))
  for i, dbFeed := range dbFeeds {
    feedIds[i] = dbFeed.ID
  }

  return jd.Refresh(feedIds)
}

func (jd *Journalistd) Refresh(feedIds []uuid.UUID) ([]error) {
  var errs []error

  dbFeeds, err := jd.entClient.Feed.
    Query().
    Where(
      feed.IDIn(feedIds...),
    ).
    All(context.Background())
  if err != nil {
    errs = append(errs, err)
    return errs
  }

  for _, dbFeed := range dbFeeds {
    rc, err := rss.NewClient(
      dbFeed.URL,
      dbFeed.Username,
      dbFeed.Password,
      true,
    )
    if err != nil {
      errs = append(errs, err)
      continue
    }

    dbFeedTmp := jd.entClient.Feed.
      Create()
    rc.SetFeed(
      dbFeed.FeedFeedLink,
      dbFeed.Username,
      dbFeed.Password,
      dbFeedTmp,
    )
    err = dbFeedTmp.
      OnConflict().
      UpdateNewValues().
      Exec(context.Background())
    if err != nil {
      errs = append(errs, err)
    }

    dbItems := make([]*ent.ItemCreate, len(rc.Feed.Items))
    for i := 0; i < len(rc.Feed.Items); i++ {
      dbItems[i] = jd.entClient.Item.
        Create()
      dbItems[i] = rc.SetItem(
        i,
        dbItems[i],
      )
    }
    err = jd.entClient.Item.
      CreateBulk(dbItems...).
      OnConflict().
      UpdateNewValues().
      Exec(context.Background())
    if err != nil {
      errs = append(errs, err)
    }
  }

  return errs
}
