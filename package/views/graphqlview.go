package views

import (
	"context"
	"encoding/json"
	"github.com/golang/glog"
	"github.com/graphql-go/graphql"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/mcodes"
	"io/ioutil"
	"net/http"
)

type GraphQLView struct {
	Context *core2.HContext
}

func (gqlView *GraphQLView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("PATH: " + r.URL.Path)

	uid := ""
	ok := false
	group := ""

	// check for auth
	if gqlView.Context.HasAuth() {
		h := httpGetHeader(r.Header, core2.HeaderAuthorizationCanonical)
		if h == "" {
			GraphQLErrorMessage(w,
				[]byte("Oops! No Authorization header found !!!"),
				mcodes.NoAuthorizationHeaderFound, 400)

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
				GraphQLErrorMessage(w,
					[]byte("Oops! No UID found from AuthVerifyInterface !!!"),
					mcodes.NoUIDFromAuthVerifyInterface, 400)
				return
			}

		} else { // Failed to authorize. not ok
			GraphQLErrorMessage(w, []byte("Oops! Invalid Authorization data !!!"),
				mcodes.InvalidAuthorizationData, 400)
			return
		}
	} // end of auth

	glog.Infoln("Processing Middleware in the GraphQLView.")
	for _, mi := range gqlView.Context.MiddlewareList {
		next := mi.ServeHTTP(gqlView.Context, r, w)
		if next == false {
			return
		}
	}
	glog.Infoln("Middleware processing complete.")

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

	if gqlView.Context.HasAuth() {
		params = graphql.Params{
			Schema:        gqlView.Context.GraphQLSchema,
			RequestString: string(bodyBytes),
			Context: context.WithValue(context.Background(),
				"auth",
				map[string]string{"uid": uid, "group": group}),
		}
	} else {
		params = graphql.Params{
			Schema:        gqlView.Context.GraphQLSchema,
			RequestString: string(bodyBytes),
			Context: context.WithValue(context.Background(),
				"auth",
				map[string]string{}),
		}
	}

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
