package views

import (
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	"github.com/mkawserm/hypcore/package/constants"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/cors"
	"github.com/mkawserm/hypcore/package/gqltypes"
	"github.com/mkawserm/hypcore/package/mcodes"
	"github.com/mkawserm/hypcore/package/variants"
	"io/ioutil"
	"net/http"
)

type GraphQLView struct {
	Context *core2.HContext
}

func (gqlView *GraphQLView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("PATH: " + r.URL.Path)

	if !cors.CheckCROSAndStepForward(gqlView.Context.CORSOptions, w, r) {
		glog.Infoln("CORS!!!")
		return
	}

	glog.Infoln("Processing Middleware in the GraphQLView")
	for _, mi := range gqlView.Context.MiddlewareList {
		next := mi.ServeHTTP(gqlView.Context, r, w)
		if next == false {
			return
		}
	}
	glog.Infoln("Middleware processing complete")

	uid := ""
	ok := false
	group := ""

	// check for auth
	if gqlView.Context.HasAuthVerify() {
		h := httpGetHeader(r.Header, constants.HeaderAuthorizationCanonical)
		if h == "" {
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.GraphQLGroupCode
			errorType.Code = mcodes.GraphQLNoAuthorizationHeaderFound
			errorType.MessageType = "GraphQLException"

			errorType.AddStringMessage("Oops! No Authorization header found !!!")

			GraphQLSmartErrorMessage(w, errorType, 400)

			return

		} else {
			var dataMap map[string]interface{}

			dataMap, ok = gqlView.Context.AuthVerify.Verify([]byte(h),
				[]byte(gqlView.Context.AuthPublicKey),
				[]byte(gqlView.Context.AuthSecretKey),
				gqlView.Context.AuthBearer)

			if ok {
				if uniqueId, found := dataMap["uid"]; found {
					uid = uniqueId.(string)
				}
				if groupString, found := dataMap["group"]; found {
					group = groupString.(string)
				}
			}
		}

		if ok {
			if uid == "" {
				errorType := gqltypes.NewErrorType()
				errorType.Group = mcodes.GraphQLGroupCode
				errorType.Code = mcodes.GraphQLNoUIDFoundFromToken
				errorType.MessageType = "GraphQLException"
				errorType.AddStringMessage("Oops! No UID found from AuthVerifyInterface !!!")
				GraphQLSmartErrorMessage(w, errorType, 400)
				return
			}

		} else { // Failed to authorize. not ok
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.GraphQLGroupCode
			errorType.Code = mcodes.GraphQLInvalidAuthorizationData
			errorType.MessageType = "GraphQLException"
			errorType.AddStringMessage("Oops! Invalid Authorization data !!!")
			GraphQLSmartErrorMessage(w, errorType, 400)
			return
		}
	} // end of auth

	if r.Method != http.MethodPost {
		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.GraphQLGroupCode
		errorType.Code = mcodes.GraphQLQueryMustBeUsingPostRequest
		errorType.MessageType = "GraphQLException"
		errorType.AddStringMessage("Oops! GraphQL query must be done using post request !!!")
		GraphQLSmartErrorMessage(w, errorType, 400)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.GraphQLGroupCode
		errorType.Code = mcodes.GraphQLRequestBodyReadError
		errorType.MessageType = "GraphQLException"
		errorType.AddStringMessage("Oops! Failed to read request body !!!")
		GraphQLSmartErrorMessage(w, errorType, 400)
		return
	}

	ro := variants.ParseGraphQLQuery(bodyBytes)
	//fmt.Println(string(bodyBytes))

	if ro == nil {
		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.GraphQLGroupCode
		errorType.Code = mcodes.GraphQLRequestBodyParseError
		errorType.MessageType = "GraphQLException"
		errorType.AddStringMessage("Oops! Failed to parse request body !!!")
		GraphQLSmartErrorMessage(w, errorType, 400)
		return
	}

	contextData := make(map[string]interface{})
	contextData["auth"] = map[string]string{"uid": uid, "group": group}

	if gqlView.Context.EnableTLS {
		contextData["connectionTpe"] = constants.HTTPSConnection
	} else {
		contextData["connectionTpe"] = constants.HTTPConnection
	}

	var params graphql.Params
	params = graphql.Params{
		Schema:         *gqlView.Context.GraphQLSchema,
		RequestString:  ro.Query,
		VariableValues: ro.Variables,
		OperationName:  ro.OperationName,

		Context: context.WithValue(context.Background(),
			"contextData", contextData),
	}

	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		glog.Errorf("failed to execute graphql operation, errors: %+v\n", res.Errors)

		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.GraphQLGroupCode
		errorType.Code = mcodes.GraphQLExecutionError
		errorType.MessageType = "GraphQLException"

		errorType.Messages = make([]interface{}, 0)
		errorType.Messages = append(errorType.Messages,
			map[string]string{"message": "Oops! GraphQL query execution error!!!"})

		for _, formattedErr := range res.Errors {
			errorType.AddMessage(formattedErr)
		}

		GraphQLSmartErrorMessage(w, errorType, 400)
		return
	}

	rJSON, _ := json.Marshal(res)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rJSON)
}
