package core

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
)

type ServeWSGraphQL struct {
}

func (serveWSGraphQL *ServeWSGraphQL) ServeWS(context *HContext, connectionId int, message []byte) {

	params := graphql.Params{Schema: context.GraphQLSchema, RequestString: string(message)}
	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		glog.Errorln("ServeWS: failed to execute graphql operation, errors: %+v", res.Errors)

		output := fmt.Sprintf("{\"message\":\"%s\", \"code\":%d}",
			[]byte("Oops! GraphQL query execution error!!!"),
			400)
		context.WriteMessage(connectionId, []byte(output))

	} else {
		rJSON, _ := json.Marshal(res)
		context.WriteMessage(connectionId, rJSON)
	}
}
