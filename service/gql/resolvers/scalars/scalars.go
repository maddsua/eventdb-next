package scalars

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func MkPtr[T any](val T) *T {
	return &val
}

func NullEpoch(timestamp *time.Time) sql.NullInt64 {

	if timestamp == nil {
		return sql.NullInt64{}
	}

	return sql.NullInt64{Int64: timestamp.Unix(), Valid: true}
}

func NullUuidString(id *uuid.UUID) sql.NullString {

	if id == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: id.String(), Valid: true}
}

func UnwrapNullString(val sql.NullString) *string {

	if !val.Valid {
		return nil
	}

	return &val.String
}

func NullString(val *string) sql.NullString {

	if val == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: *val, Valid: true}
}
