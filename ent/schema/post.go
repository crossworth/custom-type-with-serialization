package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/microcosm-cc/bluemonday"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return nil
}

// Mixin of the Post.
func (Post) Mixin() []ent.Mixin {
	return []ent.Mixin{
		SanitizeHook{
			bluemonday: bluemonday.UGCPolicy(),
		},
	}
}
