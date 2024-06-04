package resolver

import "github.com/otakakot/sample-go-gqlgen/pkg/gql/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos map[string][]*model.Todo
}

func New() *Resolver {
	todos := make(map[string][]*model.Todo)

	return &Resolver{
		todos: todos,
	}
}
