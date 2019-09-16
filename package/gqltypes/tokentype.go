package gqltypes

import "github.com/graphql-go/graphql"

var TokenType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Token",
	Description: "JSON Web Token",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "JWT encoded string",
		},
	},
})

type Token struct {
	Token string `json:"token"`
}
