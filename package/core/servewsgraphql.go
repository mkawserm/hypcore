package core

import (
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"github.com/mkawserm/hypcore/package/constants"
	"github.com/mkawserm/hypcore/package/gqltypes"
	"github.com/mkawserm/hypcore/package/mcodes"
	"github.com/mkawserm/hypcore/package/variants"
)

type ServeWSGraphQL struct {
}

func (serveWSGraphQL *ServeWSGraphQL) ServeWS(ctx *HContext, connectionId int, message []byte) {

	ro := variants.ParseGraphQLQuery(message)
	if ro == nil {
		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.GraphQLWSGroupCode
		errorType.Code = mcodes.GraphQLWSMessageParseError
		errorType.MessageType = "GraphQLWSException"
		errorType.AddStringMessage("Oops! Failed to parse request body !!!")
		ctx.WriteMessage(connectionId, GraphQLSmartErrorMessageBytes(errorType))
		return
	}

	var params graphql.Params
	contextData := make(map[string]interface{})
	contextData["auth"] = map[string]string{
		"uid":   ctx.GetUIDFromSID(connectionId),
		"group": ctx.GetGroupFromSID(connectionId)}

	if ctx.EnableTLS {
		contextData["connectionTpe"] = constants.WebSocketSecureConnection
	} else {
		contextData["connectionTpe"] = constants.WebSocketConnection
	}

	params = graphql.Params{
		Schema:        *ctx.GraphQLSchema,
		RequestString: string(message),
		Context:       context.WithValue(context.Background(), "contextData", contextData),
	}

	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		glog.Errorf("ServeWS: failed to execute graphql operation, errors: %+v\n", res.Errors)

		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.GraphQLWSGroupCode
		errorType.Code = mcodes.GraphQLWSExecutionError
		errorType.MessageType = "GraphQLWSException"
		errorType.AddStringMessage("Oops! GraphQL query execution error. Invalid query!!!")
		for _, formattedErr := range res.Errors {
			errorType.AddMessage(formattedErr)
		}
		ctx.WriteMessage(connectionId, GraphQLSmartErrorMessageBytes(errorType))

	} else {
		rJSON, _ := json.Marshal(res)
		ctx.WriteMessage(connectionId, rJSON)
	}
}
