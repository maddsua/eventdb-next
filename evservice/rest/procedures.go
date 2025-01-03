package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	sqliteops "github.com/maddsua/eventdb-next/database/operations/sqlite"
)

func NewHandler(db *sqliteops.Queries, prefix string) http.Handler {

	procs := Procedues{DB: db}

	mux := http.NewServeMux()

	mux.Handle("POST /push", handlerFn(procs.Push))

	return http.StripPrefix(strings.TrimRight(prefix, "/"), mux)
}

type Procedues struct {
	DB *sqliteops.Queries
}

func (this *Procedues) Push(writer http.ResponseWriter, req *http.Request) (any, error) {

	var checkAuth = func(streamID string, streamKey null.String) error {

		//	todo: cache check result in redis or some shi

		if entry, err := this.DB.GetStreamByID(req.Context(), streamID); err != nil {

			if this.DB.IsNoRows(err) {
				return &ProcErr{
					Message: "invalid stream id key",
					Status:  http.StatusBadRequest,
				}
			}

			return &ProcErr{
				Message: "failed to get stream",
				Status:  http.StatusInternalServerError,
			}

		} else if entry.PushKey.Valid && (!streamKey.Valid || entry.PushKey.String != streamKey.String) {
			return &ProcErr{
				Message: "invalid stream key",
				Status:  http.StatusBadRequest,
			}
		}

		return nil
	}

	var entries []modelPushEventEntry
	var streamID string

	switch req.Header.Get("content-type") {
	case "application/json":
		if data, err := ParseJSON[modelPushEventPayload](req); err != nil {
			return nil, err
		} else if err = checkAuth(data.StreamID, data.StreamKey); err != nil {
			return nil, err
		} else {
			entries = append(entries, data.modelPushEventEntry)
			streamID = data.StreamID
		}
	case "application/json+batch":
		if data, err := ParseJSON[modelPushEventBatchPayload](req); err != nil {
			return nil, err
		} else if err = checkAuth(data.StreamID, data.StreamKey); err != nil {
			return nil, err
		} else {
			entries = data.Entries
			streamID = data.StreamID
		}
	default:
		return nil, errors.New("unexpected content type")
	}

	go func() {

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		tx, err := this.DB.BeginTx(ctx)
		if err != nil {
			slog.Error("faield to start a tx",
				slog.String("err", err.Error()))
		}

		defer tx.Rollback()

		var count int

		for _, entry := range entries {

			var fields null.String
			if entry.Fields != nil {

				for key, val := range entry.Fields {
					switch val.(type) {
					case int, *int, bool, *bool:
						break
					case string, *string:
						//	todo: check length
						break
					default:
						delete(entry.Fields, key)
					}
				}

				if len(entry.Fields) > 0 {
					if data, err := json.Marshal(entry.Fields); err == nil {
						fields = null.StringFrom(string(data))
					}
				}
			}

			//	todo: check and normalize data
			arg := sqliteops.AddEventParams{
				ID:            uuid.NewString(),
				CreatedAt:     entry.Timestamp.NullInt64,
				StreamID:      streamID,
				TransactionID: entry.TransactionID.NullString,
				ClientIp:      entry.ClientIP.NullString,
				Type:          null.StringFrom("LOG").NullString,
				Level:         null.StringFrom(strings.ToLower(entry.Level)).NullString,
				Message:       null.StringFrom(entry.Message).NullString,
				Fields:        fields.NullString,
			}

			if err := tx.AddEvent(ctx, arg); err != nil {
				slog.Error("faield to insert even entry",
					slog.String("err", err.Error()))
			} else {
				count++
			}
		}

		if count == 0 {
			return
		}

		if err := tx.Commit(); err != nil {
			slog.Error("faield to commit a tx",
				slog.String("err", err.Error()))
		}
	}()

	return nil, nil
}
