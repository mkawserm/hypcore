package gqltypes

import "github.com/graphql-go/graphql"

var TokenType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Token",
	Fields: graphql.Fields{
		"token": &graphql.Field{
			Type: graphql.String,
		},
	},
})

type Token struct {
	Token string `json:"token"`
}
