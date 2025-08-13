package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Authorize holds the schema definition for the Authorize entity.
type Authorize struct {
	ent.Schema
}

// Fields of the Authorize.
func (Authorize) Fields() []ent.Field {
	return []ent.Field{
		// 注意：ent 会自动为您创建一个 'id' 字段，其类型默认为 int，
		// 在 MySQL 中会映射为 bigint。
		// 因此，我们不需要在这里显式定义 'id' 字段。
		field.String("client_id").
			Comment("客户端ID"),

		field.String("client_secret").
			Comment("客户端密钥"),

		field.String("grant_type").
			Comment("授权类型"),

		field.String("redirect_uri").
			Comment("重定向URI"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新时间"),

		field.Bool("deleted").
			Default(false).
			Comment("是否删除"),
	}
}

// Annotations of the Authorize.
func (Authorize) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "authorize"},
	}
}

// Edges of the Authorize.
func (Authorize) Edges() []ent.Edge {
	return nil
}
