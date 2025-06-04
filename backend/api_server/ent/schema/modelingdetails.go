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

// ModelingDetails holds the schema definition for the ModelingDetails entity.
type ModelingDetails struct {
	ent.Schema
}

func (ModelingDetails) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "modeling_details"},
		entsql.WithComments(true),
		schema.Comment("Modeling result table"),
	}
}

func (ModelingDetails) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("modeling_id", "model", "data_type").Unique(),
	}
}

// Fields of the ModelingDetails.
func (ModelingDetails) Fields() []ent.Field {
	defaultValue := []string{}
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.Int("modeling_id").Optional().Default(0).Comment("Parent modeling ID"),
		field.String("model").Default("").Comment("model name"),
		field.String("data_type").Default("").Comment("result data type"),
		field.JSON("data", []string{}).Default(defaultValue).Comment("result data"),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Edges of the ModelingDetails.
func (ModelingDetails) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("modeling", Modeling.Type).
			Ref("modeling_details").
			Unique().
			Field("modeling_id"),
	}
}
