package schema

import (
  "time"
  "github.com/go-playground/validator/v10"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Feed holds the schema definition for the Feed entity.
type Feed struct {
	ent.Schema
}

// Fields of the Feed.
func (Feed) Fields() []ent.Field {
  validate := validator.New()

  return []ent.Field{
    field.UUID("id", uuid.UUID{}).
      Default(uuid.New),
      // StorageKey("oid"),
    field.String("url").
      Validate(func(s string) error {
        return validate.Var(s, "required,url")
      }),
    field.String("username").
      Default("").
      Sensitive(),
    field.String("password").
      Default("").
      Sensitive(),

    field.String("feed_title"),
    field.String("feed_description"),
    field.String("feed_link"),
    field.String("feed_feed_link"),
    field.Time("feed_updated"),
    field.Time("feed_published"),
    field.String("feed_author_name").
      Optional(),
    field.String("feed_author_email").
      Optional(),
    field.String("feed_language"),
    field.String("feed_image_title").
      Optional(),
    field.String("feed_image_url").
      Optional(),
    field.String("feed_copyright"),
    field.String("feed_generator"),
    field.String("feed_categories"),

    field.Time("created_at").
      Default(time.Now),
    field.Time("updated_at").
      Default(time.Now).
      UpdateDefault(time.Now),
    field.Time("deleted_at").
      Default(nil).
      Optional().
      Nillable(),
  }
}

func (Feed) Indexes() []ent.Index {
  return []ent.Index{
    index.Fields("url", "username", "password").
      Unique(),
  }
}

// Edges of the Feed.
func (Feed) Edges() []ent.Edge {
  return []ent.Edge{
    edge.To("items", Item.Type),
    edge.From("subscribed_users", User.Type).
      Ref("subscribed_feeds").
      Through("subscriptions", Subscription.Type),
  }
}
