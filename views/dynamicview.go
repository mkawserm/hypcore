package views

import (
	"github.com/mkawserm/hypcore/core"
	"log"
	"net/http"
	"regexp"
)

type DynamicView struct {
	Context *core.HContext
}

func (dView *DynamicView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("PATH: " + r.URL.Path)

	// check for auth
	if dView.Context.HasAuth() {
		uid := ""
		ok := false

		h := httpGetHeader(r.Header, core.HeaderAuthorizationCanonical)
		if h == "" {
			httpBadRequest(w, []byte("Oops! No Authorization header found !!!"))
			return

		} else {
			uid, ok = dView.Context.Auth.GetUID([]byte(h))
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

	log.Printf("Processing Middleware in the DynamicView.")
	for _, mi := range dView.Context.MiddlewareList {
		next := mi.ServeHTTP(dView.Context, r, w)
		if next == false {
			return
		}
	}
	log.Printf("Middleware processing complete.")

	log.Printf("Dynamic route dispatch started.")
	for _, route := range dView.Context.RouteList {
		var rc = regexp.MustCompile(route.Pattern)
		if rc.Match([]byte(r.URL.Path)) {
			route.HttpHandlerObject.ServeHTTP(dView.Context, r, w)
			return
		}
	}

	log.Printf("Dynamic route dispatch failed to find a route.")
	log.Printf("Showing 404 Http Error")

	httpNotFound(w, []byte("Oops!!!"))
}
