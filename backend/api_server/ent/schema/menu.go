package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Menu holds the schema definition for the Menu entity.
type Menu struct {
	ent.Schema
}

// Annotations of the Menu.
func (Menu) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "menu"},
	}
}

// Fields of the Menu.
func (Menu) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			StorageKey("menu_key").
			StructTag(`json:"key"`),
		field.String("label"),
		field.String("icon"),
		field.String("url").
			Optional(),
		field.Bool("is_use").
			Optional().
			Default(true).
			StructTag("isUse"),
		field.Bool("is_title").
			Optional().
			Default(false).
			StructTag("isTitle"),
		field.Int("menu_order").
			StructTag("menuOrder"),
		field.String("parent_key").
			Optional().
			StructTag("parentKey"),
		field.Int("group").
			Default(2),
	}
}

// Edges of the Menu.
func (Menu) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Menu.Type).
			From("parent").
			Unique().
			Field("parent_key"),
	}
}
