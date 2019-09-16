package gqltypes

import "github.com/graphql-go/graphql"

var StatsType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Stats",
	Description: "Server statistics",
	Fields: graphql.Fields{
		"totalWebSocketConnection": &graphql.Field{
			Type:        graphql.Int,
			Description: "Total Active Web Socket Connection",
		},
	},
})

type Stats struct {
	TotalWebSocketConnection int `json:"totalWebSocketConnection"`
}
