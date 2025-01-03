package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/guregu/null"
	"github.com/maddsua/eventdb-next/utils"
)

type procedureHandler func(writer http.ResponseWriter, req *http.Request) (any, error)

func handlerFn(handlerFn procedureHandler) http.Handler {
	procHandler := procHandlerImpl{callback: handlerFn}
	return requestWrapper(procHandler)
}

type procHandlerImpl struct {
	callback procedureHandler
}

// Procedure writer acts as a proxy for standard http.ResponseWriter that allows to track
// if a procedure had written any data directly or is using function return values to return it's results
// The "Clean" propery should be set to true on writer initialization as it's only really used once
// and there's no point in writing an initializer function for it. Just remember to pass both an actual writer and a "true" to it before using
type ProcedureWriter struct {
	Writer http.ResponseWriter
	Clean  bool
}

func (this *ProcedureWriter) Header() http.Header {
	return this.Writer.Header()
}
func (this *ProcedureWriter) Write(data []byte) (int, error) {
	this.Clean = false
	return this.Writer.Write(data)
}
func (this *ProcedureWriter) WriteHeader(statusCode int) {
	this.Clean = false
	this.Writer.WriteHeader(statusCode)
}

func (this procHandlerImpl) ServeHTTP(writer http.ResponseWriter, req *http.Request) {

	procWriter := &ProcedureWriter{Writer: writer, Clean: true}

	result, err := func() (result any, err error) {

		defer func() {
			if rec := recover(); rec != nil {

				switch rect := rec.(type) {
				case error:
					err = rect
				case string:
					err = errors.New(rect)
				default:
					err = errors.New("unhandled handler exception")
				}

				slog.Error("REST: handler crashed",
					slog.String("err", err.Error()),
					slog.String("at", string(debug.Stack())),
					slog.String("rid", requestID(req)),
					slog.String("ip", req.RemoteAddr))
			}
		}()

		result, err = this.callback(procWriter, req)
		return
	}()

	//	don't serialize data if a procedure had any writes to response directly
	if !procWriter.Clean {
		return
	}

	var werr error
	if err != nil {

		status := http.StatusBadRequest
		exts := map[string]string{}
		if errx, ok := err.(*ProcErr); ok {
			status = errx.StatusCode()
			exts = errx.Extensions
		}

		payload := ResponsePayload{Error: null.StringFrom(err.Error()), ErrExt: exts}
		werr = WriteResponse(writer, payload, status)

	} else {
		payload := ResponsePayload{Data: result}
		werr = WriteResponse(writer, payload, http.StatusOK)
	}

	if werr != nil {
		slog.Error("REST: Failed to write handler response",
			slog.String("err", err.Error()),
			slog.String("ip", req.RemoteAddr),
			slog.String("rid", requestID(req)))
	}
}

type ResponsePayload struct {
	Data   any               `json:"data"`
	Error  null.String       `json:"error,omitempty"`
	ErrExt map[string]string `json:"error_extensions,omitempty"`
}

type ProcErr struct {
	Message    string
	Status     int
	Extensions map[string]string
}

func (this *ProcErr) Error() string {
	return this.Message
}

func (this *ProcErr) StatusCode() int {

	if this.Status < http.StatusBadRequest {
		return http.StatusBadRequest
	}

	return this.Status
}

func WriteResponse(writer http.ResponseWriter, resp ResponsePayload, status int) error {

	if status == http.StatusOK && resp.Data == nil {
		writer.WriteHeader(http.StatusNoContent)
		return nil
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(resp); err != nil {
		return nil
	}

	return nil
}

func ParseJSON[T any](req *http.Request) (T, error) {

	var result T

	switch req.Method {
	case http.MethodPost:
		break
	default:
		return result, errors.New("unsupported method")
	}

	if !strings.Contains(req.Header.Get("Content-Type"), "json") {
		return result, errors.New("unsupported content type")
	}

	return result, json.NewDecoder(req.Body).Decode(&result)
}

func requestID(req *http.Request) string {
	return utils.ReadRequestID(req)
}

func requestWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		req.RemoteAddr = utils.ClientIP(req)
		utils.RequestID(writer, req)
		next.ServeHTTP(writer, req)
	})
}
