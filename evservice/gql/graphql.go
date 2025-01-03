package gql

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	sqliteops "github.com/maddsua/eventdb-next/database/operations/sqlite"
	"github.com/maddsua/eventdb-next/gql/request"
	"github.com/maddsua/eventdb-next/gql/resolvers"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func NewHandler(db *sqliteops.Queries) http.Handler {

	resolver := resolvers.Resolver{DB: db}

	gqlHandler := handler.New(resolvers.NewExecutableSchema(resolvers.Config{Resolvers: &resolver}))

	gqlHandler.AddTransport(transport.SSE{})
	gqlHandler.AddTransport(transport.POST{})

	gqlHandler.Use(extension.Introspection{})

	gqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	gqlHandler.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	gqlHandler.SetRecoverFunc(func(ctx context.Context, err any) (userMessage error) {

		req := request.From(ctx)

		slog.Error("GQL: Handler crashed",
			slog.Any("err", err),
			slog.String("ip", req.ClientIP),
			slog.String("rid", req.RequestID),
			slog.String("stack", string(debug.Stack())))

		return errors.New("failed to execute query: internal server error")
	})

	gqlHandler.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {

		gqlErr := graphql.DefaultErrorPresenter(ctx, err)
		//	todo: fix
		if len(gqlErr.Extensions) == 0 {
			req := request.From(ctx)
			slog.Error("GQL: Resolver error",
				slog.Any("err", err),
				slog.String("ip", req.ClientIP),
				slog.String("rid", req.RequestID))
		}

		return gqlErr
	})

	//	todo: fix
	//gqlHandler.Use(extension.FixedComplexityLimit(env.GetIntInRange("GQL_MAX_COMPLEXITY", 10, 10_000, 200)))

	return request.InjectConext(gqlHandler)
}
