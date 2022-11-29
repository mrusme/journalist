package journalistd

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/mrusme/journalist/ent"
	"github.com/mrusme/journalist/ent/feed"
	"github.com/mrusme/journalist/ent/item"
	"github.com/mrusme/journalist/ent/user"

	"github.com/mrusme/journalist/lib"
	"github.com/mrusme/journalist/rss"
)

var VERSION string

type Journalistd struct {
	jctx *lib.JournalistContext

	config    *lib.Config
	entClient *ent.Client
	logger    *zap.Logger

	daemonStop          chan bool
	autoRefreshInterval time.Duration
}

func Version() string {
	return VERSION
}

func New(
	jctx *lib.JournalistContext,
) (*Journalistd, error) {
	jd := new(Journalistd)
	jd.jctx = jctx
	jd.config = jctx.Config
	jd.entClient = jctx.EntClient
	jd.logger = jctx.Logger

	if err := jd.initAdminUser(); err != nil {
		return nil, err
	}

	interval, err := strconv.Atoi(jd.config.Feeds.AutoRefresh)
	if err != nil {
		jd.logger.Error(
			"Feeds.AutoRefresh is not a valid number (seconds)",
			zap.Error(err),
		)
		return nil, err
	}
	jd.autoRefreshInterval = time.Duration(interval)

	return jd, nil
}

func (jd *Journalistd) IsDebug() bool {
	debug, err := strconv.ParseBool(jd.config.Debug)
	if err != nil {
		return false
	}

	return debug
}

func (jd *Journalistd) initAdminUser() error {
	var admin *ent.User
	var err error

	admin, err = jd.entClient.User.
		Query().
		Where(user.Username(jd.config.Admin.Username)).
		Only(context.Background())
	if err != nil {
		admin, err = jd.entClient.User.
			Create().
			SetUsername(jd.config.Admin.Username).
			SetPassword(jd.config.Admin.Password).
			SetRole("admin").
			Save(context.Background())
		if err != nil {
			jd.logger.Error(
				"Failed query/create admin user",
				zap.Error(err),
			)
			return err
		}
	}

	if admin.Password == "admin" {
		jd.logger.Debug(
			"Admin user",
			zap.String("username", admin.Username),
			zap.String("password", admin.Password),
		)
	} else {
		jd.logger.Debug(
			"Admin user",
			zap.String("username", admin.Username),
			zap.String("password", "xxxxxx"),
		)
	}

	return nil
}

func (jd *Journalistd) Start() bool {
	jd.logger.Info(
		"Starting Journalist daemon",
	)
	jd.daemonStop = make(chan bool)
	go jd.daemon()
	return true
}

func (jd *Journalistd) Stop() {
	jd.logger.Info(
		"Stopping Journalist daemon",
	)
	jd.daemonStop <- true
}

func (jd *Journalistd) daemon() {
	jd.logger.Debug(
		"Journalist daemon started, looping",
	)
	for {
		select {
		case <-jd.daemonStop:
			jd.logger.Debug(
				"Journalist daemon loop ended",
			)
			return
		default:
			jd.logger.Debug(
				"RefreshAll starting, refreshing all feeds",
			)
			errs := jd.RefreshAll()
			if len(errs) > 0 {
				jd.logger.Error(
					"RefreshAll completed with errors",
					zap.Errors("errors", errs),
				)
			} else {
				jd.logger.Debug(
					"RefreshAll completed",
				)
			}
			time.Sleep(time.Second * jd.autoRefreshInterval)
		}
	}
}

func (jd *Journalistd) RefreshAll() []error {
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

func (jd *Journalistd) Refresh(feedIds []uuid.UUID) []error {
	var errs []error

	dbFeeds, err := jd.entClient.Feed.
		Query().
		Where(
			feed.IDIn(feedIds...),
		).
		WithItems(func(query *ent.ItemQuery) {
			query.
				Select(item.FieldItemGUID).
				Where(item.CrawlerContentHTMLNEQ(""))
		}).
		All(context.Background())
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	for _, dbFeed := range dbFeeds {
		var exceptItemGUIDs []string
		for _, exceptItem := range dbFeed.Edges.Items {
			exceptItemGUIDs = append(exceptItemGUIDs, exceptItem.ItemGUID)
		}

		rc, errr := rss.NewClient(
			dbFeed.URL,
			dbFeed.Username,
			dbFeed.Password,
			true,
			exceptItemGUIDs,
			jd.logger,
		)
		if len(errr) > 0 {
			errs = append(errs, errr...)
			continue
		}

		dbFeedTmp := jd.entClient.Feed.
			Create()
		rc.SetFeed(
			dbFeed.URL,
			dbFeed.Username,
			dbFeed.Password,
			dbFeedTmp,
		)
		feedID, err := dbFeedTmp.
			OnConflict().
			UpdateNewValues().
			ID(context.Background())
		if err != nil {
			errs = append(errs, err)
		}

		dbItems := make([]*ent.ItemCreate, len(rc.Feed.Items))
		for i := 0; i < len(rc.Feed.Items); i++ {
			dbItems[i] = jd.entClient.Item.
				Create()
			dbItems[i] = rc.SetItem(
				feedID,
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
