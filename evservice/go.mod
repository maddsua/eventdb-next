module github.com/maddsua/eventdb-next

go 1.23.2

replace github.com/maddsua/eventdb-next/database => ../database

require (
	github.com/99designs/gqlgen v0.17.62
	github.com/google/uuid v1.6.0
	github.com/vektah/gqlparser/v2 v2.5.21
)

require (
	github.com/agnivade/levenshtein v1.2.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
)
