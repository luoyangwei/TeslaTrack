package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("account").Comment("账号"),
		field.String("password"),
		field.String("mobile").Optional(),
		field.String("open_id").Optional(),
		field.String("avatar").Optional(),
		field.String("nick_name").Optional(),
		field.String("introduction").Optional(),
		field.Int8("gender").Nillable().Default(0),
		field.Int("asked_user_id").Optional(),
		field.String("area_code").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Bool("deleted").Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "user"},
		schema.Comment("用户表"),
	}
}
