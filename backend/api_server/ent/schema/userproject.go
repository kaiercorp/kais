package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserProject holds the schema definition for the UserProject entity.
type UserProject struct {
	ent.Schema
}

// Annotations of the UserProject.
func (UserProject) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "kais_user_project"},
	}
}

// Fields of the User.
func (UserProject) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Unique().
			Immutable(),
		field.Int("project_id"),
		field.String("username"),
		field.Bool("is_use").
			Default(true),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			Optional(),
	}
}

// Edges of the UserProject.
func (UserProject) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("user_project").
			Unique().
			Field("project_id").
			Required(),
	}
}
