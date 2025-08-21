package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Partner holds the schema definition for the Partner entity.
type Partner struct {
	ent.Schema
}

// Fields of the Partner.
func (Partner) Fields() []ent.Field {
	return []ent.Field{
		field.String("client_id"),
		field.String("access_token").
			Optional(),
		field.Int("expires_in").
			Optional(),
		field.String("token_type").
			MaxLen(125).
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Bool("deleted").
			Default(false),
	}
}

// Edges of the Partner.
func (Partner) Edges() []ent.Edge {
	return nil
}

// Annotations of the Partner.
func (Partner) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "partner"},
		schema.Comment("合作伙伴"),
	}
}
