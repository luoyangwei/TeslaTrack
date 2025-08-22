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
		field.String("vin").NotEmpty().Comment("Vehicle VIN code"),
		field.Int("user_id").Comment("Associated user ID"),
		field.String("display_name").Comment("Vehicle display name"),
		field.String("access_type").Comment("Access type, e.g., OWNER"),
		field.Int8("state").Comment("Vehicle state, e.g., online, offline"),
		field.Int8("in_service").Nillable().Default(0).Comment("Is vehicle in service"),
		field.Int8("calendar_enabled").Nillable().Default(0).Comment("Is calendar enabled"),
		field.String("car_type").Optional().Comment("Car model type"),
		field.String("api_version").Optional().Comment("API version used by vehicle"),
		field.String("raw_data").Comment("Raw vehicle data from API"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("Creation time"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("Update time"),
		field.Bool("deleted").Default(false).Comment("Is deleted"),
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
		schema.Comment("Vehicle table"),
	}
}
