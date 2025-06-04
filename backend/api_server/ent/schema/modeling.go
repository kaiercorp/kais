package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Modeling holds the schema definition for the Modeling entity.
type Modeling struct {
	ent.Schema
}

func (Modeling) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "modeling"},
		entsql.WithComments(true),
		schema.Comment("Modeling task table"),
	}
}

// Fields of the Modeling.
func (Modeling) Fields() []ent.Field {
	defaultValue := []string{}
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.Int("local_id").Default(1).Comment("Modeling index in each task"),
		field.Int("task_id").Optional().Default(0).Comment("Parent task ID"),
		field.Int("parent_id").Optional().Default(0).Comment("Base Modeling ID"),
		field.Int("parent_local_id").Optional().Default(0),
		field.Int("dataset_id").Optional().Default(0).Comment("Dataset ID"),
		field.JSON("params", []string{}).Default(defaultValue).Comment("User configuration"),
		field.JSON("dataset_stat", []string{}).Default(defaultValue).Comment("Engine에서 측정한 dataset 정보"),
		field.String("modeling_type").Default("modeling").Comment("initial | update | evaluation"),
		field.String("modeling_step").Default("idle").Comment("current task step"),
		field.JSON("performance", []string{}).Default(defaultValue).Comment("Performance"),
		field.Float("progress").Comment("task progress"),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now),
		field.Time("started_at").Optional().Nillable(),
	}
}

// Edges of the Modeling.
func (Modeling) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("task", Task.Type).
			Ref("modelings").
			Unique().
			Field("task_id"),
		edge.To("modeling_details", ModelingDetails.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("modeling_models", ModelingModels.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("trials", Trial.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
