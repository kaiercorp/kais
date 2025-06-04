package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// TrialStatus holds the schema definition for the TrialStatus entity.
type TrialStatus struct {
	ent.Schema
}

func (TrialStatus) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "trial_status"},
		entsql.WithComments(true),
		schema.Comment("Trial epoch table"),
	}
}

// Fields of the TrialStatus.
func (TrialStatus) Fields() []ent.Field {
	defaultValue := []string{}
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.String("trial_uuid").Optional().Comment("Parent trial UUID"),
		field.JSON("status_json", []string{}).Default(defaultValue).Comment("epoch data"),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Edges of the TrialStatus.
func (TrialStatus) Edges() []ent.Edge {
	return nil
}
