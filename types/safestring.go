package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/microcosm-cc/bluemonday"
)

var p = bluemonday.UGCPolicy()

type SafeString string

func (s *SafeString) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		*s = SafeString(v)
	default:
		return fmt.Errorf("got unexpected type %T", v)
	}

	return nil
}

func (s SafeString) Value() (driver.Value, error) {
	return p.Sanitize(string(s)), nil
}
