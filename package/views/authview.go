package views

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/mcodes"
	"io/ioutil"
	"net/http"
)

type AuthView struct {
	Context *core2.HContext
}

func (authView *AuthView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("PATH: " + r.URL.Path)

	if r.Method != http.MethodPost {
		GraphQLErrorMessage(w,
			[]byte("Oops! GraphQL query must be done using post request !!!"),
			mcodes.InvalidRequestMethod, 400)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		GraphQLErrorMessage(w, []byte("Oops! Failed to read request body !!!"),
			mcodes.FailedToReadRequestBody, 400)
		return
	}

	var params graphql.Params

	params = graphql.Params{Schema: authView.Context.AuthSchema,
		RequestString: string(bodyBytes)}

	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		glog.Errorf("failed to execute graphql operation, errors: %+v\n", res.Errors)
		httpBadRequest(w, []byte("Oops! GraphQL query execution error. Invalid query!!!"))
		return
	}

	rJSON, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rJSON)
}
