package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("ad").MaxLen(100),
		field.String("email").Unique(),
		field.String("password").Sensitive(),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("url", Url.Type),
	}
}
