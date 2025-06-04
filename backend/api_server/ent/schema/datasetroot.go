package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DatasetRoot holds the schema definition for the DatasetRoot entity.
type DatasetRoot struct {
	ent.Schema
}

func (DatasetRoot) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "dataset_root"},
	}
}

// Fields of the DatasetRoot.
func (DatasetRoot) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.String("name").Default("Dataset"),
		field.String("path").Default("/workspace/data"),
		field.Bool("is_use").Default(true),
	}
}

// Edges of the DatasetRoot.
func (DatasetRoot) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("datasets", Dataset.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
