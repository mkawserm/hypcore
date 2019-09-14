package views

import (
	"github.com/golang/glog"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/mcodes"
	"net/http"
	"regexp"
)

type DynamicView struct {
	Context *core2.HContext
}

func (dView *DynamicView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("PATH: " + r.URL.Path)

	// check for auth
	if dView.Context.HasAuth() {
		uid := ""
		ok := false

		h := httpGetHeader(r.Header, core2.HeaderAuthorizationCanonical)
		if h == "" {
			GraphQLErrorMessage(w, []byte("Oops! No Authorization header found !!!"),
				mcodes.NoAuthorizationHeaderFound, 400)
			return

		} else {
			uid, ok = dView.Context.AuthVerify.GetUID([]byte(h), dView.Context.AuthBearer)
		}

		if ok {
			if uid == "" {
				GraphQLErrorMessage(w, []byte("Oops! No UID found from AuthVerifyInterface !!!"),
					mcodes.NoUIDFromAuthInterface, 400)
				return
			}

		} else { // Failed to authorize. not ok
			GraphQLErrorMessage(w, []byte("Oops! Invalid Authorization data !!!"),
				mcodes.InvalidAuthorizationData, 400)
			return
		}
	}

	glog.Infoln("Processing Middleware in the DynamicView.")
	for _, mi := range dView.Context.MiddlewareList {
		next := mi.ServeHTTP(dView.Context, r, w)
		if next == false {
			return
		}
	}
	glog.Infoln("Middleware processing complete.")

	glog.Infoln("Dynamic route dispatch started.")
	for _, route := range dView.Context.RouteList {
		var rc = regexp.MustCompile(route.Pattern)
		if rc.Match([]byte(r.URL.Path)) {
			route.HttpHandlerObject.ServeHTTP(dView.Context, r, w)
			return
		}
	}

	glog.Infoln("Dynamic route dispatch failed to find a route.")
	glog.Errorln("Showing 404 Http Error")

	GraphQLErrorMessage(w, []byte("Oops!!!"), mcodes.HttpNotFound, 404)
}
