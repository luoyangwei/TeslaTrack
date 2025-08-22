package schema

import "entgo.io/ent"

// Vehicle holds the schema definition for the Vehicle entity.
type Vehicle struct {
	ent.Schema
}

// Fields of the Vehicle.
func (Vehicle) Fields() []ent.Field {
	return nil
}

// Edges of the Vehicle.
func (Vehicle) Edges() []ent.Edge {
	return nil
}
