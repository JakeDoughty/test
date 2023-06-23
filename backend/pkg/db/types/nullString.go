package types

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type NullString struct {
	String string
	Valid  bool
}

func (s NullString) ToString() string {
	if !s.Valid {
		return ""
	} else {
		return s.String
	}
}
func (s NullString) GoString() string {
	if s.Valid {
		return fmt.Sprintf("%q", s.String)
	} else {
		return "nil"
	}
}
func (s *NullString) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*s = NullString{}
		return nil

	case string:
		*s = NullString{String: v, Valid: true}
		return nil

	default:
		return errors.New("text string")
	}
}
func (s NullString) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}

	return s.Value, nil
}
func (s NullString) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "TEXT"
}
func (s NullString) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if !s.Valid {
		return gorm.Expr("?", nil)
	} else {
		return gorm.Expr("?", s.String)
	}
}

func ValidString(s string) NullString {
	return NullString{String: s, Valid: true}
}
func StringToNullString(s string) NullString {
	if s == "" {
		return NullString{}
	} else {
		return ValidString(s)
	}
}
func StringPtrToNullString(s *string) NullString {
	if s == nil {
		return NullString{}
	} else {
		return ValidString(*s)
	}
}
