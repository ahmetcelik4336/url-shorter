package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Url struct {
	ent.Schema
}

func (Url) Fields() []ent.Field {
	return []ent.Field{
		field.String("short_code").MaxLen(255).NotEmpty().Unique(),
		field.Text("long_url").NotEmpty(),
		field.Time("created_at").Immutable().Default(time.Now),
		field.String("alias").Optional(),
		field.Bool("is_encrypted").Optional().Default(false),
		field.String("password").Optional().Sensitive(),
		field.Time("expiration_date").Optional(),
	}
}

func (Url) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("url").
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.SetNull,
			}),
		edge.To("log_url", Logs.Type),
	}
}
