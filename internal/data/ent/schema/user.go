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
		field.String("account").Comment("Account"),
		field.String("password").Comment("Password"),
		field.String("mobile").Optional().Comment("Mobile phone number"),
		field.String("open_id").Optional().Comment("Wechat OpenID"),
		field.String("avatar").Optional().Comment("User avatar URL"),
		field.String("nick_name").Optional().Comment("User nickname"),
		field.String("introduction").Optional().Comment("User introduction"),
		field.Int8("gender").Nillable().Default(0).Comment("Gender, 0:unknown, 1:male, 2:female"),
		field.Int("asked_user_id").Optional().Comment("ID of the user who invited this user"),
		field.String("area_code").Optional().Comment("Area code for mobile number"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Creation time"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Update time"),
		field.Bool("deleted").Default(false).Comment("Is deleted"),
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
		schema.Comment("User table"),
	}
}
