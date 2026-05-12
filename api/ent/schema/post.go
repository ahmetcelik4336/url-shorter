package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Post struct {
	ent.Schema
}

func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("post_title").MaxLen(255),
		field.Text("post_content").Optional(),
	}
}

func (Post) Edges() []ent.Edge {
	return nil

}
