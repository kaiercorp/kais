package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Trial holds the schema definition for the Trial entity.
type Trial struct {
	ent.Schema
}

func (Trial) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "trial"},
		entsql.WithComments(true),
		schema.Comment("trials in modeling"),
	}
}

// Fields of the Trial.
func (Trial) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.Int("modeling_id").Optional().Default(0),
		field.String("uuid"),
		field.String("state"),
		field.String("save_path"),
		field.String("target_metric").Default("wa"),
		field.Float("progress").Default(0.0),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).Optional(),
	}
}

// Edges of the Trial.
func (Trial) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("modeling", Modeling.Type).
			Ref("trials").
			Unique().
			Field("modeling_id"),
	}
}
