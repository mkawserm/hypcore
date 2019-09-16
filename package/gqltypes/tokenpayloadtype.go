package gqltypes

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

func parseLiteral(astValue ast.Value) interface{} {
	kind := astValue.GetKind()

	switch kind {
	case kinds.StringValue:
		return astValue.GetValue()
	case kinds.BooleanValue:
		return astValue.GetValue()
	case kinds.IntValue:
		return astValue.GetValue()
	case kinds.FloatValue:
		return astValue.GetValue()
	case kinds.ObjectValue:
		obj := make(map[string]interface{})
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			obj[v.Name.Value] = parseLiteral(v.Value)
		}
		return obj
	case kinds.ListValue:
		list := make([]interface{}, 0)
		for _, v := range astValue.GetValue().([]ast.Value) {
			list = append(list, parseLiteral(v))
		}
		return list
	default:
		return nil
	}
}

// JSON json type
var JSON = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "JSON",
		Description: "The `JSON` scalar type represents JSON values",
		Serialize: func(value interface{}) interface{} {
			return value
		},
		ParseValue: func(value interface{}) interface{} {
			return value
		},
		ParseLiteral: parseLiteral,
	},
)

var TokenPayloadType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "TokenPayload",
	Description: "JSON Web Token Payload",
	Fields: graphql.Fields{
		"isValid": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Is JSON Web Token valid",
		},
		"payload": &graphql.Field{
			Type:        JSON,
			Description: "JSON Web Token key value pair",
		},
	},
})

type TokenPayload struct {
	IsValid bool                   `json:"isValid"`
	Payload map[string]interface{} `json:"payload"`
}
