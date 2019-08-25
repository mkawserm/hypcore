package views

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/mkawserm/hypcore/core"
	"io/ioutil"
	"log"
	"net/http"
)

type GraphQLView struct {
	Context *core.HContext
}

func (gqlView *GraphQLView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("PATH: " + r.URL.Path)

	// check for auth
	if gqlView.Context.HasAuth() {
		uid := ""
		ok := false

		h := httpGetHeader(r.Header, core.HeaderAuthorizationCanonical)
		if h == "" {
			httpBadRequest(w, []byte("Oops! No Authorization header found !!!"))
			return

		} else {
			uid, ok = gqlView.Context.Auth.GetUID([]byte(h))
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

	log.Printf("Processing Middleware in the GraphQLView.")
	for _, mi := range gqlView.Context.MiddlewareList {
		next := mi.ServeHTTP(gqlView.Context, r, w)
		if next == false {
			return
		}
	}
	log.Printf("Middleware processing complete.")

	if r.Method != http.MethodPost {
		httpBadRequest(w, []byte("Oops! GraphQL query must be done using post request !!!"))
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpBadRequest(w, []byte("Oops! Failed to read request body !!!"))
		return
	}

	params := graphql.Params{Schema: gqlView.Context.GraphQLSchema, RequestString: string(bodyBytes)}
	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		log.Printf("failed to execute graphql operation, errors: %+v", res.Errors)
		httpBadRequest(w, []byte("Oops! GraphQL query execution error!!!"))
		return
	}

	rJSON, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(rJSON)
}
