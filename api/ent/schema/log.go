package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Logs struct {
	ent.Schema
}

// $env:GOWORK="off"; go generate ./ent
func (Logs) Fields() []ent.Field {
	return []ent.Field{
		field.String("device"),
		field.String("ip"),
		field.String("referer"),
		field.Time("created_at"),
		field.String("type"),
	}
}

func (Logs) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("log", Url.Type).
			Ref("log_url").
			Unique().
			Required().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
