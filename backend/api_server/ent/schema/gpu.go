package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Gpu holds the schema definition for the Gpu entity.
type Gpu struct {
	ent.Schema
}

func (Gpu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "gpu"},
	}
}

// Fields of the Gpu.
func (Gpu) Fields() []ent.Field {
	return []ent.Field{
		field.String("uuid").Unique(),                // GPU UUID
		field.Int("index").Default(0),                // GPU identifier in machine
		field.String("name"),                         // GPU product name
		field.String("state").Default("idle"),        // idle | train | test | load
		field.Bool("is_use").Default(true),           // usable
		field.Int("device_id").Optional().Default(0), // GPU machine ID
	}
}

// Edges of the Gpu.
func (Gpu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("device", Device.Type).
			Ref("gpu").
			Unique().
			Field("device_id"),
	}
}
