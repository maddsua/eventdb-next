package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	sqliteops "github.com/maddsua/eventdb-next/database/operations/sqlite"
	"github.com/maddsua/eventdb-next/gql/resolvers/scalars"
)

func TransformStreamEvent(val sqliteops.Event) (StreamEvent, error) {

	id, err := uuid.Parse(val.ID)
	if err != nil {
		return StreamEvent{}, err
	}

	streamID, err := uuid.Parse(val.StreamID)
	if err != nil {
		return StreamEvent{}, err
	}

	return StreamEvent{
		ID:   id,
		Date: time.Unix(val.CreatedAt, 0),
		Stream: DataStream{
			ID: streamID,
		},
		ClientIP:      scalars.UnwrapNullString(val.ClientIp),
		TransactionID: scalars.UnwrapNullString(val.TransactionID),
		Level:         transformLogLevel(val.Level.String),
		Message:       val.Message.String,
		Fields:        transformEventFields(val.Fields),
	}, nil
}

func transformLogLevel(val string) LogLevel {

	level := LogLevel(val)
	if !level.IsValid() {
		return LogLevelError
	}

	return level
}

func transformEventFields(val sql.NullString) []EventField {

	if !val.Valid {
		return nil
	}

	var fields map[string]any
	_ = json.Unmarshal([]byte(val.String), &fields)

	var entries []EventField
	for key, val := range fields {
		entries = append(entries, EventField{
			Key:   key,
			Value: fmt.Sprintf("%v", val),
		})
	}

	return entries
}

func TransformDataStream(val sqliteops.Stream) (DataStream, error) {

	id, err := uuid.Parse(val.ID)
	if err != nil {
		return DataStream{}, err
	}

	return DataStream{
		ID:      id,
		Created: time.Unix(val.CreatedAt, 0),
		Updated: time.Unix(val.UpdatedAt, 0),
		Name:    val.Name,
		PushKey: scalars.UnwrapNullString(val.PushKey),
	}, nil
}
