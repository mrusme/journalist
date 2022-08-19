package schema

import (
  "time"
  "github.com/go-playground/validator/v10"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
    field.String("title"),
    field.String("description"),
    field.String("site_url").
      Validate(func(s string) error {
        return validate.Var(s, "required,url")
      }),
    field.String("feed_url").
      Validate(func(s string) error {
        return validate.Var(s, "required,url")
      }),
    field.String("author"),
    field.String("language"),
    field.String("image"),
    field.String("copyright"),
    field.String("generator"),
    field.String("categories"),
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

// Edges of the Feed.
func (Feed) Edges() []ent.Edge {
  return []ent.Edge{
    edge.To("items", Item.Type),
    edge.From("subscribed_users", User.Type).
      Ref("subscribed_feeds").
      Through("subscriptions", Subscription.Type),
  }
}
