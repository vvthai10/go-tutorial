//go:generate granate
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/vvthai10/demo-graphql/bootstrap"
	"github.com/vvthai10/demo-graphql/graphql"
	"github.com/vvthai10/demo-graphql/repository/postgres"
)

func main() {
	app := bootstrap.App()
	userRepo := postgres.NewUserRepository(app.DB)

	ctx := context.Background()

	var gqls graphql.Service
	{
		schema, err := graphql.NewSchema(
			graphql.NewResolver(userRepo),
		)
		if err != nil {
			os.Exit(1)
		}
		gqls = graphql.NewService(schema)
	}

	http.Handle("/graphql", graphql.MakeHandler(ctx, gqls))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}
