package schema

import (
  "time"
  "github.com/go-playground/validator/v10"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/edge"
	"github.com/google/uuid"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
  validate := validator.New()

  return []ent.Field{
    field.UUID("id", uuid.UUID{}).
      Default(uuid.New),
      // StorageKey("oid"),

    field.String("item_title"),
    field.String("item_description"),
    field.String("item_content"),
    field.String("item_link").
      Validate(func(s string) error {
        return validate.Var(s, "required,url")
      }),
    field.String("item_updated"),
    field.String("item_published"),
    field.String("item_author"),
    field.String("item_authors"),
    field.String("item_guid"),
    field.String("item_image"),
    field.String("item_categories"),
    field.String("item_enclosures"),

    field.String("crawler_title"),
    field.String("crawler_author"),
    field.String("crawler_excerpt"),
    field.String("crawler_site_name"),
    field.String("crawler_image"),
    field.String("crawler_content_html"),
    field.String("crawler_content_text"),

    field.Time("created_at").
      Default(time.Now),
    field.Time("updated_at").
      Default(time.Now).
      UpdateDefault(time.Now),
  }
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
  return []ent.Edge{
    edge.From("feed", Feed.Type).
      Ref("items").
      Unique(),
    edge.From("read_by_users", User.Type).
      Ref("read_items").
      Through("reads", Read.Type),
  }
}
