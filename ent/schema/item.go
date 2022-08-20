package schema

import (
  "time"
  "github.com/go-playground/validator/v10"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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

    field.String("item_guid").
      Unique(),
    field.String("item_title"),
    field.String("item_description"),
    field.String("item_content"),
    field.String("item_link").
      Validate(func(s string) error {
        return validate.Var(s, "required,url")
      }),
    field.String("item_updated"),
    field.String("item_published"),
    field.String("item_author_name").
      Optional(),
    field.String("item_author_email").
      Optional(),
    field.String("item_image_title").
      Optional(),
    field.String("item_image_url").
      Optional(),
    field.String("item_categories"),
    field.String("item_enclosures"),

    field.String("crawler_title").
      Optional(),
    field.String("crawler_author").
      Optional(),
    field.String("crawler_excerpt").
      Optional(),
    field.String("crawler_site_name").
      Optional(),
    field.String("crawler_image").
      Optional(),
    field.String("crawler_content_html").
      Optional(),
    field.String("crawler_content_text").
      Optional(),

    field.Time("created_at").
      Default(time.Now),
    field.Time("updated_at").
      Default(time.Now).
      UpdateDefault(time.Now),
  }
}

func (Item) Indexes() []ent.Index {
  return []ent.Index{
    index.Fields("item_guid").
      Unique(),
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
