package schema

import (
	"entgo.io/bug/types"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age"),
		field.String("name"),
		field.Other("content", types.SafeString("")).SchemaType(map[string]string{
			dialect.Postgres: "varchar",
			dialect.MySQL:    "varchar",
			dialect.SQLite:   "varchar",
		}),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
