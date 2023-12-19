package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/vvthai10/demo-graphql/entities"
	"github.com/vvthai10/demo-graphql/service/user"
	"golang.org/x/net/context"
)

type resolver struct {
	userService user.IUserService
}

func NewResolver(userRepo user.IUserRepository) *resolver {
	return &resolver{
		userService: user.NewService(userRepo),
	}
}

func (r *resolver) getUser(ctx context.Context, email string) (*entities.User, error) {
	return r.userService.GetUserByEmail(ctx, email)
}

func NewSchema(r *resolver) (graphql.Schema, error) {
	userType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"fullName": &graphql.Field{
					Type: graphql.String,
				},
				"email": &graphql.Field{
					Type: graphql.String,
				},
				"password": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	fields := graphql.Fields{
		"user": &graphql.Field{
			Description: "User",
			Type:        userType,
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if email, ok := p.Args["email"].(string); ok {
					if user, err := r.getUser(p.Context, email); err != nil {
						return nil, err
					} else {
						return user, nil
					}
				}
				return nil, nil
			},
		},
	}

	rootQueryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: fields,
		},
	)
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQueryType,
	})
}
