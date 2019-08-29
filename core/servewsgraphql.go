package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
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

		output := fmt.Sprintf("{\"message\":\"%s\", \"error_code\":%d}",
			[]byte("Oops! GraphQL query execution error. Invalid query!!!"),
			400)
		ctx.WriteMessage(connectionId, []byte(output))

	} else {
		rJSON, _ := json.Marshal(res)
		ctx.WriteMessage(connectionId, rJSON)
	}
}
