package views

import (
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	core2 "github.com/mkawserm/hypcore/package/core"
	"io/ioutil"
	"net/http"
)

type GraphQLView struct {
	Context *core2.HContext
}

func (gqlView *GraphQLView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("PATH: " + r.URL.Path)

	uid := ""

	// check for auth
	if gqlView.Context.HasAuth() {
		ok := false

		h := httpGetHeader(r.Header, core2.HeaderAuthorizationCanonical)
		if h == "" {
			httpBadRequest(w, []byte("Oops! No Authorization header found !!!"))
			return

		} else {
			uid, ok = gqlView.Context.Auth.GetUID([]byte(h), gqlView.Context.AuthBearer)
		}

		if ok {
			if uid == "" {
				httpBadRequest(w, []byte("Oops! No UID found from AuthInterface !!!"))
				return
			}

		} else { // Failed to authorize. not ok
			httpBadRequest(w, []byte("Oops! Invalid Authorization data !!!"))
			return
		}
	}

	glog.Infoln("Processing Middleware in the GraphQLView.")
	for _, mi := range gqlView.Context.MiddlewareList {
		next := mi.ServeHTTP(gqlView.Context, r, w)
		if next == false {
			return
		}
	}
	glog.Infoln("Middleware processing complete.")

	if r.Method != http.MethodPost {
		httpBadRequest(w, []byte("Oops! GraphQL query must be done using post request !!!"))
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpBadRequest(w, []byte("Oops! Failed to read request body !!!"))
		return
	}

	var params graphql.Params

	if gqlView.Context.HasAuth() {
		params = graphql.Params{
			Schema:        gqlView.Context.GraphQLSchema,
			RequestString: string(bodyBytes),
			Context: context.WithValue(context.Background(),
				"auth",
				map[string]string{"uid": uid}),
		}
	} else {
		params = graphql.Params{Schema: gqlView.Context.GraphQLSchema, RequestString: string(bodyBytes)}
	}

	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		glog.Errorln("failed to execute graphql operation, errors: %+v", res.Errors)
		httpBadRequest(w, []byte("Oops! GraphQL query execution error. Invalid query!!!"))
		return
	}

	rJSON, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rJSON)
}
