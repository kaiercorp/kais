package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Configuration holds the schema definition for the Configuration entity.
type Configuration struct {
	ent.Schema
}

// Annotations of the Configuration.
func (Configuration) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "config"},
	}
}

func (Configuration) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("config_type", "config_key").Unique(),
	}
}

// Fields of the Configuration.
func (Configuration) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.String("config_type"),
		field.String("config_key"),
		field.String("config_val"),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("updated_at").Default(time.Now).Optional(),
	}
}

// Edges of the Configuration.
func (Configuration) Edges() []ent.Edge {
	return nil
}
