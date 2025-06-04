package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// EngineLog holds the schema definition for the EngineLog entity.
type EngineLog struct {
	ent.Schema
}

func (EngineLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "enginelog"},
	}
}

// Fields of the EngineLog.
func (EngineLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.Int("modeling_id").Default(-1),
		field.String("level"),
		field.String("filename"),
		field.Int("line"),
		field.String("message"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the EngineLog.
func (EngineLog) Edges() []ent.Edge {
	return nil
}
