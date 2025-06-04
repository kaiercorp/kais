package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// HyperParamsHistory holds the schema definition for the HyperParamsHistory entity.
type HyperParamsHistory struct {
	ent.Schema
}

func (HyperParamsHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "hyper_params_history"},
		entsql.WithComments(true),
		schema.Comment("Hyper parameter history per trial"),
	}
}

// Fields of the HyperParamsHistory.
func (HyperParamsHistory) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.Int("trial_id").Default(0).Comment("modeling_id"),
		field.String("trial_uuid").Default(""),
		field.String("model").Default(""),
		field.Int("model_num").Default(0),
		field.String("params").Default(""),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Edges of the HyperParamsHistory.
func (HyperParamsHistory) Edges() []ent.Edge {
	return nil
}
