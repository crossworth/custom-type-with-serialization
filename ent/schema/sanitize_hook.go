package schema

import (
	"context"

	"entgo.io/bug/ent/hook"
	"entgo.io/ent"
	"entgo.io/ent/schema/mixin"
	"github.com/microcosm-cc/bluemonday"
)

type SanitizeHook struct {
	mixin.Schema
	bluemonday *bluemonday.Policy
}

func (s SanitizeHook) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				for _, fieldName := range m.Fields() {
					fVal, _ := m.Field(fieldName)

					switch v := fVal.(type) {
					case string:
						if err := m.SetField(fieldName, s.bluemonday.Sanitize(v)); err != nil {
							return nil, err
						}
					}
				}

				return next.Mutate(ctx, m)
			})
		}, ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),
	}
}
