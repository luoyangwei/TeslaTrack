package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Vehicle holds the schema definition for the Vehicle entity.
type Vehicle struct {
	ent.Schema
}

// Fields of the Vehicle.
func (Vehicle) Fields() []ent.Field {
	return []ent.Field{
		field.String("vin").NotEmpty().Comment("车辆 VIN 码"),
		field.Int("user_id"),
		field.String("display_name"),
		field.String("access_type"),
		field.Int8("state"),
		field.Int8("in_service"),
		field.Int8("calendar_enabled").Optional(),
		field.String("car_type").Optional(),
		field.String("api_version").Optional(),
		field.String("raw_data").Comment("原始数据"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
		field.Bool("deleted").Default(false),
	}
}

// Edges of the Vehicle.
func (Vehicle) Edges() []ent.Edge {
	return nil
}

// Annotations of the Vehicle.
func (Vehicle) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "vehicle"},
		schema.Comment("车辆"),
	}
}
