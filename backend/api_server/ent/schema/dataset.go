package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Dataset holds the schema definition for the Dataset entity.
type Dataset struct {
	ent.Schema
}

func (Dataset) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "dataset"},
	}
}

// Fields of the Dataset.
func (Dataset) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Int("parent_id").Optional(),
		field.String("description").Optional(),
		field.String("path"),
		field.Bool("is_valid").Default(false),
		field.Bool("is_trainable").Default(false),
		field.Bool("is_testable").Default(false),
		field.Bool("is_leaf").Default(false),
		field.Bool("is_deleted").Default(false),
		field.Bool("is_use").Default(true),
		field.JSON("stat", []string{}),
		field.String("stat_path").Optional(),
		field.Strings("engine").Optional(),
		field.String("data_type"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
		field.Time("deleted_at").Default(time.Now),
		field.Int("dr_id").Optional().Default(0),
	}
}

// Edges of the Dataset.
func (Dataset) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("datasetroot", DatasetRoot.Type).
			Ref("datasets").
			Unique(),
	}
}
