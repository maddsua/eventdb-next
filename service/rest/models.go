package rest

import "github.com/guregu/null"

type modelPushEventPayload struct {
	modelPushEventEntry
	StreamID  string      `json:"stream_id"`
	StreamKey null.String `json:"stream_key"`
}

type modelPushEventBatchPayload struct {
	StreamID  string                `json:"stream_id"`
	StreamKey null.String           `json:"stream_key"`
	Entries   []modelPushEventEntry `json:"entries"`
}

type modelPushEventEntry struct {
	modulePushEventMetadata
	Message string         `json:"message"`
	Fields  map[string]any `json:"fields,omitempty"`
}

type modulePushEventMetadata struct {
	Timestamp     null.Int    `json:"timestamp"`
	Level         string      `json:"level"`
	ClientIP      null.String `json:"client_ip"`
	TransactionID null.String `json:"transaction_id"`
}
