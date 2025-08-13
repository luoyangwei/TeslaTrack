package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// AuthorizeToken holds the schema definition for the AuthorizeToken entity.
// This entity stores the OAuth 2.0 tokens required to interact with the Tesla API.
type AuthorizeToken struct {
	ent.Schema
}

// Fields of the AuthorizeToken.
func (AuthorizeToken) Fields() []ent.Field {
	return []ent.Field{
		// The authorization code received from Tesla's OAuth 2.0 flow.
		field.String("tesla_code"),
		// The client ID for your application registered with Tesla.
		field.String("client_id"),
		// The client secret for your application registered with Tesla.
		field.String("client_secret"),
		// The access token used to make authenticated requests to the Tesla API.
		field.String("access_token"),
		// The refresh token used to obtain a new access token when the current one expires.
		field.String("refresh_token"),
		// The scope of permissions granted by the access token (e.g., "vehicle_data").
		field.String("scope"),
		// The time the token record was created.
		field.Time("created_at").Default(time.Now),
		// The time the token record was last updated.
		field.Time("updated_at").Default(time.Now),
		// A flag indicating whether the token record has been soft-deleted.
		field.Bool("deleted").Default(false),
	}
}

// Annotations of the AuthorizeToken.
func (AuthorizeToken) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Sets the table name in the database to "authorize_token".
		entsql.Annotation{Table: "authorize_token"},
	}
}

// Edges of the AuthorizeToken.
// An edge defines a relationship to another entity.
func (AuthorizeToken) Edges() []ent.Edge {
	// This entity currently has no defined edges (relationships) to other entities.
	return nil
}
