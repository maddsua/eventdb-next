package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/maddsua/eventdb-next/database"
	sqliteops "github.com/maddsua/eventdb-next/database/operations/sqlite"
	"github.com/maddsua/eventdb-next/gql"
	"github.com/maddsua/eventdb-next/rest"
	_ "github.com/mattn/go-sqlite3"
)

const endpointRest = "/rest/v2/"
const endpointGraphql = "/graphql"
const endpointGraphqlPg = endpointGraphql + "/playground"
const sqliteFileName = "eventdb.db3"

func main() {

	const serverPort = 8080

	var sqliteDbExists bool
	if _, err := os.Stat(sqliteFileName); err == nil {
		sqliteDbExists = true
	}

	db, err := sql.Open("sqlite3", sqliteFileName+"?_fk=true&_journal=WAL")
	if err != nil {
		slog.Error("error opening sqlite db",
			slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	if !sqliteDbExists {

		slog.Info("STARTUP: Applying sqlite schema...")

		if _, err := db.Exec(database.DbSchemaSqlite); err != nil {
			slog.Error("failed to set up sqlite schema",
				slog.String("err", err.Error()))
			os.Exit(1)
		}
	}

	dbq := sqliteops.New(db)

	rootMux := http.NewServeMux()
	rootMux.Handle(endpointRest, rest.NewHandler(dbq, endpointRest))
	rootMux.Handle(endpointGraphql, gql.NewHandler(dbq))

	rootMux.Handle(endpointGraphqlPg, playground.Handler("GraphQL playground", endpointGraphql))
	slog.Info("STARTUP: GraphQL playground enabled",
		slog.String("url", fmt.Sprintf("http://localhost:%d%s", serverPort, endpointGraphqlPg)))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: rootMux,
	}

	errorSignal := make(chan error)
	exitSignal := make(chan os.Signal, 1)

	signal.Notify(exitSignal, os.Interrupt, os.Kill)

	slog.Info("STARTUP: HTTP server listening on",
		slog.String("url", fmt.Sprintf("http://localhost:%d", serverPort)))

	go func() {

		err := server.ListenAndServe()
		if err == nil {
			return
		}

		select {
		case errorSignal <- err:
		default:
		}
	}()

	slog.Info("SERVICE: --> Ready")

	select {
	case err := <-errorSignal:
		slog.Error("SERVICE: Server crashed!",
			slog.String("err", err.Error()))
		return
	case <-exitSignal:
		slog.Warn("SERVICE: Stopping...")
		return
	}
}
