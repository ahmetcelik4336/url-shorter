package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Setting struct {
	ent.Schema
}

// $env:GOWORK="off"; go generate ./ent
func (Setting) Fields() []ent.Field {
	return []ent.Field{
		field.String("settings_key"),
		field.Text("setting_content"),
	}
}

func (Setting) Edges() []ent.Edge {
	return []ent.Edge{}
}
