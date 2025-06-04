package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

func (Task) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "task"},
		entsql.WithComments(true),
		schema.Comment("KAI.S가 제공하는 task: modeling|update model|performance evaluation"),
	}
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	defaultValue := []string{}
	return []ent.Field{
		field.Int("id").Unique().Immutable().Comment(""),
		field.Int("project_id").Optional().Default(0).Comment("Project ID"),
		field.Int("dataset_id").Comment("dataset id"),
		field.String("title").Default("Modeling task").Comment("task name"),
		field.String("description").Default("").Comment("user memo"),
		field.String("engine_type").Comment("engine type"),
		field.String("target_metric").Default("wa").Comment(""),
		field.JSON("params", []string{}).Default(defaultValue).Comment("params by engine type"),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project", Project.Type).
			Ref("tasks").
			Unique().
			Field("project_id"),
		edge.To("modelings", Modeling.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
