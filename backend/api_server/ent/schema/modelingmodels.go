package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ModelingModels holds the schema definition for the ModelingModels entity.
type ModelingModels struct {
	ent.Schema
}

func (ModelingModels) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "modeling_models"},
		entsql.WithComments(true),
		schema.Comment("model rank table"),
	}
}

func (ModelingModels) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("modeling_id", "data_type").Unique(),
	}
}

// Fields of the ModelingModels.
func (ModelingModels) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.Int("modeling_id").Optional().Default(0).Comment("Parent modeling ID"),
		field.String("data_type").Default("").Comment("result data type"),
		field.String("data").Default("{}").Comment("result data"),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Edges of the ModelingModels.
func (ModelingModels) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("modeling", Modeling.Type).
			Ref("modeling_models").
			Unique().
			Field("modeling_id"),
	}
}
