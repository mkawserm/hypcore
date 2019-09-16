package views

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/cors"
	"github.com/mkawserm/hypcore/package/gqltypes"
	"github.com/mkawserm/hypcore/package/mcodes"
	"github.com/mkawserm/hypcore/package/variants"
	"io/ioutil"
	"net/http"
)

type AuthView struct {
	Context *core2.HContext
}

func (authView *AuthView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("PATH: " + r.URL.Path)

	if !cors.CheckCROSAndStepForward(authView.Context.CORSOptions, w, r) {
		glog.Infoln("CORS!!!")
		return
	}

	if r.Method != http.MethodPost {
		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.AuthGroupCode

		errorType.Code = mcodes.AuthQueryMustBeUsingPostRequest
		errorType.MessageType = "GraphQueryException"
		errorType.AddStringMessage("Oops! GraphQL query must be done using post request !!!")

		GraphQLSmartErrorMessage(w, errorType, 400)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.AuthGroupCode
		errorType.Code = mcodes.AuthRequestBodyReadError
		errorType.MessageType = "GraphQueryException"
		errorType.AddStringMessage("Oops! Failed to read request body !!!")
		GraphQLSmartErrorMessage(w, errorType, 400)
		return
	}

	ro := variants.ParseGraphQLQuery(bodyBytes)
	//fmt.Println(string(bodyBytes))

	if ro == nil {
		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.AuthGroupCode
		errorType.Code = mcodes.AuthRequestBodyParseError
		errorType.MessageType = "GraphQueryException"
		errorType.AddStringMessage("Oops! Failed to parse request body !!!")

		GraphQLSmartErrorMessage(w, errorType, 400)
		return
	}

	//glog.Infoln(ro)
	var params graphql.Params

	params = graphql.Params{
		Schema:         *authView.Context.AuthSchema,
		RequestString:  ro.Query,
		VariableValues: ro.Variables,
		OperationName:  ro.OperationName,
	}

	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		glog.Errorf("failed to execute graphql operation, errors: %+v\n", res.Errors)

		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.AuthGroupCode
		errorType.Code = mcodes.AuthGraphQLExecutionError
		errorType.MessageType = "GraphQueryException"
		errorType.AddStringMessage("Oops! GraphQL query execution error!!!")
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
