package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TrialDetails holds the schema definition for the TrialDetails entity.
type TrialDetails struct {
	ent.Schema
}

func (TrialDetails) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "trial_details"},
		entsql.WithComments(true),
		schema.Comment("trial result table"),
	}
}

func (TrialDetails) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("trial_uuid", "model", "data_type").Unique(),
	}
}

// Fields of the TrialDetails.
func (TrialDetails) Fields() []ent.Field {
	defaultValue := []string{}
	return []ent.Field{
		field.Int("id").Unique().Immutable(),
		field.String("trial_uuid").Optional().Comment("Parent trial UUID"),
		field.String("model").Default("").Comment("model name"),
		field.String("data_type").Default("").Comment("result data type"),
		field.JSON("data", []string{}).Default(defaultValue).Comment("result data"),
		field.Bool("is_model_saved").Default(true),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

// Edges of the TrialDetails.
func (TrialDetails) Edges() []ent.Edge {
	return nil
}
