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
    field.String("title"),
    field.String("description"),
    field.String("content"),
    field.String("url").
      Validate(func(s string) error {
        return validate.Var(s, "required,url")
      }),
    field.String("author"),
    field.String("image"),
    field.String("categories"),

    field.String("crawled_title"),
    field.String("crawled_author"),
    field.String("crawled_excerpt"),
    field.String("crawled_site_name"),
    field.String("crawled_image"),
    field.String("crawled_content_html"),
    field.String("crawled_content_text"),

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
