package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"github.com/mkawserm/hypcore/package/mcodes"
)

type ServeWSGraphQL struct {
}

func (serveWSGraphQL *ServeWSGraphQL) ServeWS(ctx *HContext, connectionId int, message []byte) {
	var params graphql.Params

	if ctx.HasAuth() {
		params = graphql.Params{
			Schema:        ctx.GraphQLSchema,
			RequestString: string(message),
			Context: context.WithValue(context.Background(),
				"auth",
				map[string]string{"uid": ctx.GetUIDFromSID(connectionId)}),
		}
	} else {
		params = graphql.Params{
			Schema:        ctx.GraphQLSchema,
			RequestString: string(message),
		}
	}

	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		glog.Errorln("ServeWS: failed to execute graphql operation, errors: %+v", res.Errors)

		messageFormat := `
		{
			"error": {
						"message": "%s",
						"code": "%s"
					 }
		}`

		output := fmt.Sprintf(messageFormat, []byte("Oops! GraphQL query execution error. Invalid query!!!"),
			mcodes.InvalidGraphQLQuery)
		ctx.WriteMessage(connectionId, []byte(output))

	} else {
		rJSON, _ := json.Marshal(res)
		ctx.WriteMessage(connectionId, rJSON)
	}
}
